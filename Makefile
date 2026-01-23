
install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0 && \
	go install go.uber.org/mock/mockgen@v0.6.0 && \
	go install gotest.tools/gotestsum@v1.13.0 && \
	go mod tidy && \
	go mod vendor

test:
	go test ./...

run:
	set -a && \
	source .env && \
	set +a && \
	go run cmd/main.go

lint:
	golangci-lint run

build:
	go build -o transaction-manager cmd/main.go