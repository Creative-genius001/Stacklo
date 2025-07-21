# Stage 1 — build
FROM golang:1.24 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN go build -o main .

# Stage 2 — run
FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY .env .env

EXPOSE 8080

CMD ["./main"]
