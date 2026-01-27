import json
import logging
import random

from locust import HttpUser, between, events
from locust.exception import StopUser
from admin_utils import (
    admin_login,
    start_registration,
    pause_registration,
    reset_enrollments,
)

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)

# Configuration
ADMIN_ID = "admin"
ADMIN_PW = "1234"

# 학생당 신청할 강의 수
MIN_COURSES_TO_ENROLL = 1
MAX_COURSES_TO_ENROLL = 5

# 중복 클릭 시뮬레이션: 처음 신청할 때 몇 번 연타할지
MIN_INITIAL_CLICKS = 1
MAX_INITIAL_CLICKS = 5

students = []
courses = []
TOTAL_COURSES = 0


def load_test_data():
    global students, courses, TOTAL_COURSES

    with open("/mnt/locust/students.json", "r", encoding="utf-8") as f:
        students = json.load(f)

    with open("/mnt/locust/courses.json", "r", encoding="utf-8") as f:
        courses = json.load(f)

    TOTAL_COURSES = len(courses)
    logging.info(f"Loaded {len(students)} students and {TOTAL_COURSES} courses")


@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    load_test_data()

    # Admin login and start registration
    host = environment.host
    if host:
        logging.info(f"Starting registration via admin API at {host}")
        try:
            session = admin_login(host, ADMIN_ID, ADMIN_PW)
            # reset_enrollments는 pause 상태에서만 실행 가능
            try:
                pause_registration(host, session)
            except Exception as e:
                logging.info(f"Pause failed (may already be paused): {e}")
            reset_enrollments(host, session)
            start_registration(host, session)
        except Exception as e:
            logging.error(f"Failed to setup test: {e}")
            if environment.runner:
                environment.runner.quit()
            raise


class Student(HttpUser):
    wait_time = between(0.1, 0.5)

    def on_start(self):
        """순차적 시나리오: 로그인 → (강의조회 → 신청) 반복 → 종료"""
        # 1. 로그인
        if not self.login():
            raise StopUser()

        # 2. 신청할 강의 선택
        enrollment_plan = self.create_enrollment_plan()

        # 3. 각 강의에 대해: 강의 조회 → 신청 (중복 클릭 포함)
        for plan in enrollment_plan:
            self.wait()
            self.get_courses()

            course_id = plan["course_id"]
            clicks = plan["clicks"]

            for _ in range(clicks):
                self.wait()
                should_retry = self.enroll_course(course_id)
                if not should_retry:
                    break

        # 4. 모든 신청 완료 후 종료
        raise StopUser()

    def login(self) -> bool:
        if not students:
            logging.error("No more students available for login")
            return False

        student = students.pop()
        max_retries = 5

        for attempt in range(max_retries):
            with self.client.post(
                "/api/v1/auth/login",
                json={
                    "username": student["phone_number"],
                    "password": student["birth_date"],
                },
                catch_response=True,
            ) as response:
                if response.status_code == 200:
                    response.success()
                    self.student_info = student
                    return True
                elif response.status_code in (500, 502, 503, 504):
                    response.failure(f"HTTP {response.status_code}")
                    logging.warning(
                        f"Login attempt {attempt + 1}/{max_retries} failed: {response.status_code}"
                    )
                else:
                    response.failure(f"HTTP {response.status_code}")
                    logging.error(f"Login failed: {response.status_code}")
                    return False

        logging.error(f"Login failed after {max_retries} attempts")
        return False

    def create_enrollment_plan(self) -> list:
        """신청할 강의와 각 강의별 클릭 횟수 설정"""
        num_courses = random.randint(MIN_COURSES_TO_ENROLL, MAX_COURSES_TO_ENROLL)
        selected_course_ids = random.sample(range(1, TOTAL_COURSES + 1), num_courses)

        return [
            {
                "course_id": course_id,
                "clicks": random.randint(MIN_INITIAL_CLICKS, MAX_INITIAL_CLICKS),
            }
            for course_id in selected_course_ids
        ]

    def get_courses(self):
        self.client.get("/api/v1/courses/", name="/api/v1/courses/")

    def enroll_course(self, course_id: int) -> bool:
        """
        강의 신청 시도.
        Returns: True면 계속 클릭 시도, False면 다음 강의로
        """
        with self.client.post(
            "/api/v1/course-reg/enrollment",
            json={"course_id": course_id},
            name="/api/v1/course-reg/enrollment",
            catch_response=True,
        ) as response:
            if response.status_code == 200:
                response.success()
                return False  # 성공, 다음 강의로
            elif response.status_code == 409:
                response.success()
                return True  # 중복, 계속 클릭
            elif response.status_code in (404, 403):
                response.success()
                return False  # 실패, 다음 강의로
            else:
                response.failure(f"HTTP {response.status_code}")
                return True  # 서버 에러, 재시도
