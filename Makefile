
install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0 && \
	go install go.uber.org/mock/mockgen@v0.6.0 && \
	go install gotest.tools/gotestsum@v1.13.0 && \
	go mod tidy && \
	go mod vendor

integration-test:
	gotestsum --format pkgname ./test/integration/...

run:
	set -a && \
	source .env && \
	set +a && \
	go run cmd/main.go

lint:
	golangci-lint run

build:
	go build -o accounts-api cmd/main.go

build-docker:
	docker build -t tiagovaldrich/accounts-api .

run-docker:
	docker run -p 8889:8889 \
		-e DB_HOST=host.docker.internal:5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=postgres \
		-e DB_NAME=accounts_api \
		tiagovaldrich/accounts-api