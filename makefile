run: build
	@./bin/fableflow

build:
	@cd cmd/api; go build -o ../../bin/fableflow

dev:
	@go run ./cmd/api/main.go

test:
	@go test ./...