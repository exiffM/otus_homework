# Sender's Dockerfile
FROM golang:1.20 as build

ENV BIN_FILE /opt/sender/sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sender/*

# На выходе тонкий образ alpine 3.9 <---
FROM ubuntu:latest

RUN apt-get update
RUN apt-get install netcat -y

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/sender/sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sender/sendercfg.yml
COPY ./configs/sendercfg.yml ${CONFIG_FILE}

COPY ./wait-for.sh /opt/sender

# ENTRYPOINT ["tail" "-f" "/dev/null"]
# CMD ${BIN_FILE} -config ${CONFIG_FILE}