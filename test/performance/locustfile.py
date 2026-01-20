import json
import logging
import random

from locust import HttpUser, task, between, events
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
            pause_registration(host, session)
            reset_enrollments(host, session)
            start_registration(host, session)
        except Exception as e:
            raise Exception(f"Failed to setup test: {e}")


class Student(HttpUser):
    wait_time = between(0.1, 0.5)

    def on_start(self):
        self.student_info = self.login()

    def login(self):
        if not students:
            logging.error("No more students available for login")
            raise StopIteration("Student pool exhausted")

        student = students.pop()
        logging.info(
            f"Remaining students: {len(students)}, Current: {student['phone_number']}"
        )

        response = self.client.post(
            "/api/v1/auth/login",
            json={
                "username": student["phone_number"],
                "password": student["birth_date"],
            },
        )

        if response.status_code == 200:
            logging.info(f"Login successful: {student['phone_number']}")
        else:
            logging.error(
                f"Login failed: {response.status_code}, student: {student['phone_number']}"
            )

        return student

    @task(1)
    def get_courses(self):
        self.client.get("/api/v1/courses/", name="/api/v1/courses/")

    @task(10)
    def enroll_course(self):
        course_id = random.randint(1, TOTAL_COURSES)

        with self.client.post(
            "/api/v1/course-reg/enrollment",
            json={"course_id": course_id},
            name="/api/v1/course-reg/enrollment",
            catch_response=True,
        ) as response:
            result = response.json()
            message = result.get("message", "")

            if response.status_code == 200:
                response.success()
                logging.info(f"Enrollment success: course_id={course_id}")
            elif response.status_code in (404, 409, 403):
                # Business logic rejection (not found, conflict, forbidden)
                response.success()
                logging.info(
                    f"Enrollment rejected: course_id={course_id}, status={response.status_code}, message={message}"
                )
            else:
                response.failure(f"HTTP {response.status_code}")
                logging.error(
                    f"Enrollment request failed: course_id={course_id}, status={response.status_code}, message={message}"
                )
