BINC := "./bin/calendar"
BINSCH := "./bin/scheduler"
BINSND := "./bin/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BINC) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BINSCH) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BINSND) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) --config=./configs/config.yml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

generate:
	rm -rf ./internal/server/grpc/pb
	mkdir ./internal/server/grpc/pb

	protoc --go-grpc_out=./internal/server/grpc/pb --go_out=./internal/server/grpc/pb ./api/EventService.proto

test:
	go test -p 1 -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint -v run ./...


.PHONY: build run build-img run-img version test lint
