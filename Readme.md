# Course Registration Backend

```
docker build -t course-reg .     
docker run -v ./logs:/logs -v ./db:/db -p 3000:3000 course-reg
```

# Load Test 
```
# cd test/performance
python3 generate_test_data.py --num_students 100 --num_courses 100 --num_students 1000
python3 register_test_data.py # --reset  
docker compose -f locust-docker-compose.yml up 
```