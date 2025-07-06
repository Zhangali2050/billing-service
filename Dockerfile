# Dockerfile
FROM golang:1.24-alpine


WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./wait-for.sh", "db:5432", "--", "./main"]

