FROM golang:latest

LABEL maintainer="Sergei Kolesnikov <sergei@kolesnikov.se>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV DB_URI "host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"

RUN go build -o main .

CMD ["./main"]