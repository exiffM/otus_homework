FROM golang:1.20 AS build

# copy only goose binary
FROM ubuntu:latest

RUN apt-get update
RUN apt-get install net-tools
ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
# ENV DBSTRING postgres://igor:igorpwd@postgresqldb:5432/calendardb?sslmode=disable

COPY ./migrations/files/001_create_table.sql /app/migrations/001_create_table.sql
COPY ./migrations/inserting/002_insert_rows.sql /app/migrations/002_insert_rows.sql
# ENTRYPOINT ["/bin/bash", "-c", goose, "postgres", ${DBSTRING}, "up"]

WORKDIR /app/migrations
# ENTRYPOINT ["tail", "-f", "/dev/null"]
CMD goose up