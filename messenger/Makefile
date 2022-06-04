build:
	CGO_ENABLED=0 go build .

proto:
	protoc --go_out=./internal/delivery/proto --go-grpc_out=./internal/delivery/proto config/users.proto
	go mod tidy

lint:
	golangci-lint run
