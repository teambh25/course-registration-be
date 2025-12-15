# Course Registration Backend

```
docker build -t course-reg .     
docker run -v ./logs:/logs -v ./db:/db -p 3000:3000 course-reg
```