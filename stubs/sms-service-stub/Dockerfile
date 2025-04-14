FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG GOARCH

RUN GOOS=linux GOARCH=${GOARCH} go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 5023

CMD ["./main"]
