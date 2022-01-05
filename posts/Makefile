.DEFAULT_GOAL = build
.PHONY: test

build:
	CGO_ENABLED=0 go build -o ./posts ./cmd/posts.go

test:
	go test -count=1 ./...

test-integration:
	go test -tags=integration -count=1 ./...

proto:
	protoc --go_out=./internal/delivery/proto --go-grpc_out=./internal/delivery/proto config/posts.proto
	go mod tidy

lint:
	golangci-lint run
