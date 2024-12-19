run:
	go run cmd/api/main.go

test:
	go install github.com/vektra/mockery/v2@v2.50.0
	go generate ./...
	go test ./...