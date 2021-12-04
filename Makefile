.DEFAULT_GOAL := build

BINARY_NAME=flusso

build:
	go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test -race ./...

test_coverage:
	go test -race ./... -coverprofile=coverage.out

dep:
	go mod download

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run

format:
	goimports -w .

compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.
