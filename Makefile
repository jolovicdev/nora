build:
	@go build -o nora cmd/nora/main.go

run: build
	@go run cmd/nora/main.go