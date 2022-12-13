build:
	go build -o eh -v entrypoints/cli/main.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

ready: lint test build