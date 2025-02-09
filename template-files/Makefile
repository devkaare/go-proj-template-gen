MAIN_FILE_PATH = cmd/api/main.go

all: build test

run:
	@go run $(MAIN_FILE_PATH)

build:
	@go build $(MAIN_FILE_PATH)

test:
	@go test ./... -v

clean:
	@rm -rf main
	@go mod tidy

.PHONY: all run build test clean
