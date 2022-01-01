ARG TARGETARCH
ARG TARGETOS

FROM golang:1.17.5-alpine3.14 AS modules
ADD . /app
WORKDIR /app
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go mod download

FROM modules AS tester
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go test -v ./...

FROM modules AS builder
WORKDIR /app/cmd
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go build -o gateway .

FROM alpine:3.14
COPY --from=builder /app/cmd/gateway ./gateway
ENTRYPOINT ["./gateway"]
