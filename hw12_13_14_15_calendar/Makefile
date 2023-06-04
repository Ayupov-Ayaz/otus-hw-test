BIN := "./bin/calendar"
SENDER_BIN := "./bin/sender"
SHEDULER_BIN := "./bin/scheduler"
CONFIG := "./configs/config.yaml"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-sheduler:
	go build -v -o $(SHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sheduler

build-sender:
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run-calendar: build-calendar
	$(BIN) -config $(CONFIG)

run-sender: build-sender
	$(SENDER_BIN) -config $(CONFIG)

run-sheduler: build-sheduler
	$(SHEDULER_BIN) -config $(CONFIG)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

user ?= user
pwd ?= pwd
host ?= localhost
port ?= 3306
db ?= calendar

dsn = "mysql://${user}:${pwd}@tcp(${host}:${port})/${db}"

migrate-up:
	migrate -path ./migrations -database ${dsn} up

migrate-down:
	migrate -path ./migrations -database ${dsn} down

integration-test:
	TEST_USER=${user} TEST_PWD=${pwd} go test --tags=integration ./...

.PHONY: generate
generate:
	go generate ./...

docker-run-mysql:
	docker run --name calendar-mysql -e MYSQL_ROOT_PASSWORD=${pwd} -e MYSQL_DATABASE=${db} -p ${port}:3306 -d mysql:latest

docker-run-rabbitmq:
	docker run -d --hostname my-rabbit --name calendar-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management