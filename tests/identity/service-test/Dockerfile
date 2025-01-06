FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENTRYPOINT ["go", "test", "-v", "."]