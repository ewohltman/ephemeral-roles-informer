.PHONY: lint test build docker push deploy all

MAKEFILE_PATH=$(shell readlink -f "${0}")
MAKEFILE_DIR=$(shell dirname "${MAKEFILE_PATH}")

parentImage=alpine:latest

lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=coverage.out ./...

build:
	CGO_ENABLED=0 go build -o build/package/dbl-updater/dbl-updater cmd/dbl-updater/dbl-updater.go

image:
	docker pull "${parentImage}"
	docker image build -t ewohltman/dbl-updater:latest build/package/dbl-updater

push:
	docker login -u "${DOCKER_USER}" -p "${DOCKER_PASS}"
	docker push ewohltman/dbl-updater:latest

deploy:
	${MAKEFILE_DIR}/scripts/deploy.sh

all: lint test build image push deploy
