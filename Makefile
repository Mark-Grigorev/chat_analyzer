VERSION ?= 1.0.0
APP     := app
IMAGE   := chat_analyzer

SHELL   := /bin/bash
export PATH := $(PATH)

.PHONY: build test lint docker-build clean

build:
	go build -o $(APP) cmd/*.go

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	golangci-lint run --timeout=5m

docker-build:
	docker build -f docker/Dockerfile --build-arg VERSION=$(VERSION) -t $(IMAGE) .

clean:
	rm -f $(APP) coverage.txt
