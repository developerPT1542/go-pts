build:
	@go build -o bin/mdes-cs-server cmd/main.go

testing:
	@go test -v ./...

run: build
	@./bin/mdes-cs-server