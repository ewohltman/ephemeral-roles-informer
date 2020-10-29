.PHONY: lint test build docker push deploy all

MAKEFILE_PATH=$(shell readlink -f "${0}")
MAKEFILE_DIR=$(shell dirname "${MAKEFILE_PATH}")

version=$(shell grep 'image: ewohltman/ephemeral-roles-informer:' deployments/kubernetes/deployment.yml | awk -F: '{print $$3}')

parentImage=alpine:latest

lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=coverage.out ./...

build:
	CGO_ENABLED=0 go build -o build/package/ephemeral-roles-informer/ephemeral-roles-informer cmd/ephemeral-roles-informer/ephemeral-roles-informer.go

image:
	docker pull "${parentImage}"
	docker image build -t ewohltman/ephemeral-roles-informer:${version} build/package/ephemeral-roles-informer

push:
	docker login -u "${DOCKER_USER}" -p "${DOCKER_PASS}"
	docker push ewohltman/ephemeral-roles-informer:${version}
	docker tag ewohltman/ephemeral-roles-informer:${version} ewohltman/ephemeral-roles-informer:latest
	docker logout

deploy:
	${MAKEFILE_DIR}/scripts/deploy.sh

all: lint test build image push deploy
