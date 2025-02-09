MAIN_FILE_PATH = cmd/api/main.go

all: build test

run:
	@go run $(MAIN_FILE_PATH)

build:
	@echo "Building..."
	@go build $(MAIN_FILE_PATH)

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -rf main
	@go mod tidy

.PHONY: all run build test clean