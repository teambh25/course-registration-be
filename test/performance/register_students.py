import os
import json
import requests
import time

def send_bulk_request():
    with open('dummy_students.json', 'r', encoding='utf-8') as f:
        students = json.load(f)

    url = "http://localhost:3000"  # Endpoint for bulk insertion

    # login
    session = requests.Session()

    try:
        response = session.post(os.path.join(url, "api/v1/auth/login"), json={"username": "admin", "password": "1234"})
        print(f"Successfully login. Status code: {response.status_code}")
        print("Response:", response.json())
        print(session.cookies.get_dict())
    except requests.exceptions.RequestException as e:
        print(f"Failed to login. Error: {e}")

    try:
        cookies = session.cookies.get_dict() 
        cookie_header = "; ".join([f"{k}={v}" for k, v in cookies.items()])

        st = time.time()
        response = session.post(os.path.join(url, "api/v1/admin/students/register"), json=students, headers={"Cookie": cookie_header})
        t = time.time() - st
        response.raise_for_status()  # Raise an exception for bad status codes
        print(f"Successfully sent bulk data. Status code: {response.status_code}")
        print("Response:", response.json())
        print(f"time : {t}")
    except requests.exceptions.RequestException as e:
        print(f"Failed to send bulk data. Error: {e}")
                            
if __name__ == "__main__":
    send_bulk_request()
