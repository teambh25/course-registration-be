import os
import requests


def admin_login(url: str, admin_id: str, admin_pw: str) -> requests.Session:
    """Admin login and return authenticated session"""
    session = requests.Session()

    try:
        response = session.post(
            os.path.join(url, "api/v1/auth/login"),
            json={"username": admin_id, "password": admin_pw},
        )
        response.raise_for_status()
        print(f"Admin login successful. Status code: {response.status_code}")
        return session
    except requests.exceptions.RequestException as e:
        print(f"Admin login failed. Error: {e}")
        raise


def start_registration(url: str, session: requests.Session):
    """Start registration via admin API"""
    try:
        cookies = session.cookies.get_dict()
        cookie_header = "; ".join([f"{k}={v}" for k, v in cookies.items()])
        response = session.post(
            os.path.join(url, "api/v1/admin/registration/start"),
            headers={"Cookie": cookie_header},
        )
        response.raise_for_status()
        print(f"Registration started successfully. Status code: {response.status_code}")
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Failed to start registration. Error: {e}")
        raise
