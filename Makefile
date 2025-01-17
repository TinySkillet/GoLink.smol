DEFAULT_TARGET=build

.PHONY = fmt vet build

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@go build -o bin/smol.exe

run:
	@./bin/smol