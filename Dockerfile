FROM golang:1.25-alpine AS builder
WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# SQLite 빌드에 필요한 패키지 설치
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY . ./
RUN go mod download
RUN go build -o main .

FROM alpine:latest

# 런타임에 필요한 패키지 설치
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app
COPY --from=builder /build/main .
COPY --from=builder /build/conf/app.ini ./conf/app.ini

# DB 디렉토리 생성
RUN mkdir -p /app/db

VOLUME ["/app/db"]
EXPOSE 3000

ENTRYPOINT ["/main"]
