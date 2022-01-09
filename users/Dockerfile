FROM golang:1.17.5-alpine3.14 AS modules
ADD . /app
WORKDIR /app
RUN go mod download

FROM modules AS tester
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go test -v ./...

FROM modules AS builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /app/cmd
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go build -o users .

FROM alpine:3.14
COPY --from=builder /app/cmd/users ./users
ENTRYPOINT ["./users"]
