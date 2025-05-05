BINARY_NAME=whisker
BINARY_PATH=./bin/$(BINARY_NAME)

fmt: ## Format the code
	go fmt ./...

test: ## Run all tests
	go test ./...

run:
	go run cmd/whisker/main.go

build:
	go build -o $(BINARY_PATH) -v cmd/whisker/main.go
