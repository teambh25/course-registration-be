import os
import json
import requests
import argparse

from admin_utils import admin_login, pause_registration


def register_studnet_data(url: str, session: requests.Session):
    with open("students.json", "r", encoding="utf-8") as f:
        students = json.load(f)

    try:
        cookies = session.cookies.get_dict()
        cookie_header = "; ".join([f"{k}={v}" for k, v in cookies.items()])
        response = session.post(
            os.path.join(url, "api/v1/admin/setup/students/register"),
            json=students,
            headers={"Cookie": cookie_header},
        )
        response.raise_for_status()  # Raise an exception for bad status codes
        print(
            f"Successfully register studnet data. Status code: {response.status_code}"
        )
    except requests.exceptions.RequestException:
        print(f"Failed to registered course data : {response.json()}")


def register_course_data(url: str, session: requests.Session):
    with open("courses.json", "r", encoding="utf-8") as f:
        courses = json.load(f)

    try:
        cookies = session.cookies.get_dict()
        cookie_header = "; ".join([f"{k}={v}" for k, v in cookies.items()])
        response = session.post(
            os.path.join(url, "api/v1/admin/setup/courses/register"),
            json=courses,
            headers={"Cookie": cookie_header},
        )
        response.raise_for_status()  # Raise an exception for bad status codes
        print(f"Successfully register course data. Status code: {response.status_code}")
    except requests.exceptions.RequestException:
        print(f"Failed to registered course data : {response.json()}")


def reset_data(url: str, session: requests.Session):
    try:
        cookies = session.cookies.get_dict()
        cookie_header = "; ".join([f"{k}={v}" for k, v in cookies.items()])
        response = session.delete(
            os.path.join(url, "api/v1/admin/setup/students/reset"),
            headers={"Cookie": cookie_header},
        )
        response.raise_for_status()
        response = session.delete(
            os.path.join(url, "api/v1/admin/setup/courses/reset"),
            headers={"Cookie": cookie_header},
        )
        response.raise_for_status()
        print("Successfully reset data")
    except requests.exceptions.RequestException:
        print(f"Failed to reset data. Error: {response.json()}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--url", type=str, default="http://localhost:3000")
    parser.add_argument("--admin_id", type=str, default="admin")
    parser.add_argument("--admin_pw", type=str, default="1234")
    parser.add_argument(
        "--reset", action="store_true", help="reset before register data"
    )
    args = parser.parse_args()

    session = admin_login(args.url, args.admin_id, args.admin_pw)

    try:
        pause_registration(args.url, session)
    except Exception as e:
        print(e)

    if args.reset:
        reset_data(args.url, session)
    register_studnet_data(args.url, session)
    register_course_data(args.url, session)
