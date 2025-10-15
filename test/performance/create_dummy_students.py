import json
from faker import Faker
import random

# Initialize Faker
fake = Faker('ko_KR')

# Generate 3000 unique phone numbers
phone_numbers = set()
while len(phone_numbers) < 3000:
    phone_numbers.add(fake.phone_number())

phone_numbers = list(phone_numbers)

# Generate dummy data
dummy_students = []
for i in range(3000):
    name = fake.name()
    birth_date = fake.date_of_birth(minimum_age=20, maximum_age=30).strftime('%Y-%m-%d')
    
    dummy_students.append({
        "name": name,
        "phone_number": phone_numbers[i],
        "birth_date": birth_date
    })

# Write to a JSON file
with open('dummy_students.json', 'w', encoding='utf-8') as f:
    json.dump(dummy_students, f, ensure_ascii=False, indent=2)

print("Successfully created dummy_students.json with 3000 entries.")
