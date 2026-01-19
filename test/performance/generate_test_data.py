import random
import argparse
import json

from faker import Faker


def generate_student_data(fake: Faker, num_student: int):
    phone_numbers = [fake.unique.phone_number() for i in range(num_student)]

    # Generate dummy data
    students = []
    for i in range(num_student):
        name = fake.name()
        birth_date = fake.date_of_birth(minimum_age=1, maximum_age=99).strftime(
            "%Y-%m-%d"
        )

        students.append(
            {"name": name, "phone_number": phone_numbers[i], "birth_date": birth_date}
        )

    # Write to a JSON file
    with open("students.json", "w", encoding="utf-8") as f:
        json.dump(students, f, ensure_ascii=False, indent=2)

    print(f"Successfully created students.json with {num_student} entries.")


def generate_course_data(fake: Faker, num_course: int):
    fake = Faker("ko_KR")

    DAYS = ["월", "화", "수", "목", "금", "토", "일"]

    courses = []

    for i in range(1, num_course + 1):
        name = f"강의 {i}"
        day = random.choice(DAYS)
        start_hour = random.randint(9, 18)
        end_hour = min(start_hour + 3, 22)

        num_schedules = random.randint(1, 3)
        selected_days = random.sample(
            DAYS, k=num_schedules
        )  # 같은 요일이 중복되지 않도록 요일 풀에서 비복원 추출
        schedule_list = []
        for day in selected_days:
            start_hour = random.randint(9, 19)  # 시작 시간: 9시 ~ 19시
            duration = random.randint(1, 3)  # 수업 시간: 1 ~ 3시간
            end_hour = start_hour + duration
            time_str = f"{day} {start_hour:02d}:00-{end_hour:02d}:00"  # 포맷: "요일 HH:MM-HH:MM"
            schedule_list.append(time_str)
        schedule_list.sort(key=lambda x: DAYS.index(x.split()[0]))
        schedules_str = ", ".join(schedule_list)

        item = {
            "name": name,
            "instructor": fake.name(),
            "description": "",
            "schedules": schedules_str,
            "capacity": random.randint(1, 100),
            "is_special": random.choice([True, False]),
        }
        courses.append(item)

    with open("courses.json", "w", encoding="utf-8") as f:
        json.dump(courses, f, ensure_ascii=False, indent=2)

    print(f"Successfully created courses.json with {num_course} entries.")

    return courses


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--num_students", type=int, required=True)
    parser.add_argument("--num_courses", type=int, required=True)
    parser.add_argument("--seed", type=int, default=3)
    args = parser.parse_args()

    random.seed(args.seed)
    Faker.seed(args.seed)
    fake = Faker("ko_KR")

    generate_student_data(fake, args.num_students)
    generate_course_data(fake, args.num_courses)
