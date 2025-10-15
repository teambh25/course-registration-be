FROM golang:1.25-alpine AS builder
WORKDIR /build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

WORKDIR /app

FROM scratch
COPY --from=builder /build/main .
COPY --from=builder /build/conf/app.ini ./conf/app.ini
ENTRYPOINT ["/main"]