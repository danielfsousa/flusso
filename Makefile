.DEFAULT_GOAL := build
.PHONY: all init gencert build build-docker run build-and-run clean test test-coverage dep staticcheck lint format compile

BINARY_NAME=flusso
CONFIG_PATH=${HOME}/.flusso
DATA_PATH=${CONFIG_PATH}/data
TAG ?= 0.1.0

all: init gencert build test

init:
	mkdir -p ${DATA_PATH}

gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		test/client-csr.json | cfssljson -bare client

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="root" \
		test/client-csr.json | cfssljson -bare root-client

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="nobody" \
		test/client-csr.json | cfssljson -bare nobody-client

	mv *.pem *.csr ${CONFIG_PATH}

$(CONFIG_PATH)/model.conf:
	cp test/model.conf $(CONFIG_PATH)/model.conf

$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv $(CONFIG_PATH)/policy.csv

build:
	go build -o ${BINARY_NAME} cmd/flusso/main.go

build-docker:
	docker build -t github.com/danielfsousa/flusso:$(TAG) .

run:
	./${BINARY_NAME}

build-and-run: build run

clean:
	go clean
	rm ${BINARY_NAME}

test: $(CONFIG_PATH)/model.conf $(CONFIG_PATH)/policy.csv
	go test -race ./...

test-cov:
	go test -race ./... -coverprofile=coverage.out

dep:
	go mod download

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run

format:
	goimports -w internal test

compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.
