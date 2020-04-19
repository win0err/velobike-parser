FROM debian:latest as wait-for-it
RUN apt-get update && apt-get install wait-for-it -y

FROM golang:latest
LABEL maintainer="Sergei Kolesnikov <sergei@kolesnikov.se>"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY --from=wait-for-it /usr/bin/wait-for-it /usr/bin/wait-for-it

ENV DB_DIALECT sqlite3
ENV DB_URI "/data/velobike.db"
ENV BACKUP_DIR "/data"

RUN go build -o velobike-parser .

CMD ["./velobike-parser"]