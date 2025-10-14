# Course Registration Backend

```
# setup
go mod init course-reg
go mod tidy
docker build -t course-reg .     
docker run -p 3000:3000 course-reg
```


docker run -v logs:/logs -v db:/db -p 3000:3000 course-reg 