BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_CALENDAR_IMG="calendar:develop"
DOCKER_SCHEDULER_IMG="scheduler:develop"
DOCKER_SENDER_IMG="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run-calendar: build-calendar
	$(BIN_CALENDAR) --config=./configs/config.yml

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) --config=./configs/schedulercfg.yml

run-sender: build-sender
	$(BIN_SENDER) --config=./configs/sendercfg.yml

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_IMG) \
		-f build/calendar/Dockerfile .

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SCHEDULER_IMG) \
		-f build/scheduler/Dockerfile .

build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SENDER_IMG) \
		-f build/sender/Dockerfile .

run-calendar-img: build-calendar-img
	docker run $(DOCKER_CALENDAR_IMG)

run-scheduler-img: build-scheduler-img
	docker run $(DOCKER_SCHEDULER_IMG)

run-sender-img: build-sender-img
	docker run $(DOCKER_SENDER_IMG)

version: build
	$(BIN) version

generate:
	rm -rf ./internal/server/grpc/pb
	mkdir ./internal/server/grpc/pb

	protoc --go-grpc_out=./internal/server/grpc/pb --go_out=./internal/server/grpc/pb ./api/EventService.proto

test:
	go test -p 1 -race ./internal/...

integration-test:
	go test -tags integration -v -timeout=1m15s ./tests

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

up:
	docker compose -f deployments/docker-compose.yaml up --build

down:
	docker compose -f deployments/docker-compose.yaml down --remove-orphans

lint: install-lint-deps
	golangci-lint -v run ./...

.PHONY: build run build-img run-img version test lint
