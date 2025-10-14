FROM golang:1.25-alpine AS builder
WORKDIR /build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# COPY go.mod go.sum ./
COPY . ./
RUN go mod download


RUN go build -o main .

# COPY main.go ./
# RUN go build -o main .

WORKDIR /app
RUN cp /build/main .

FROM scratch
COPY --from=builder /app/main .
COPY --from=builder /build/conf/app.ini ./conf/app.ini
ENTRYPOINT ["/main"]
