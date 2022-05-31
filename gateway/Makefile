.DEFAULT_GOAL = build
.PHONY: test

build:
	CGO_ENABLED=0 go build .

test:
	go test -count=1 ./...

proto:
	protoc --go_out=./internal/delivery/proto --go-grpc_out=./internal/delivery/proto config/users.proto
	protoc --go_out=./internal/delivery/proto --go-grpc_out=./internal/delivery/proto config/posts.proto
	go mod tidy

lint:
	golangci-lint run --exclude logger.Sync
