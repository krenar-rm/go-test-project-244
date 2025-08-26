.PHONY: build test lint run clean

build:
	@mkdir -p bin
	go build -o ./bin/gendiff ./cmd/gendiff/main.go

test:
	go test -v ./...

test-with-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

run: build
	./bin/gendiff $(ARGS)

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
