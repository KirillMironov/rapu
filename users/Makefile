.DEFAULT_GOAL = build
.PHONY: test

build:
	CGO_ENABLED=0 go build -o ./users ./cmd/users.go

test:
	go test -count=1 ./...

test-integration:
	go test -tags=integration -count=1 ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/delivery/proto/users.proto
	go mod tidy

lint:
	golangci-lint run

