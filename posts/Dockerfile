FROM golang:1.17.5-alpine3.14 AS modules
WORKDIR /opt
COPY go.mod go.sum ./
RUN go mod download
COPY posts posts

FROM modules AS tester
WORKDIR /opt/posts
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go test -v ./...

FROM modules AS builder
WORKDIR /opt/posts/cmd
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go build -o posts .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/posts/cmd/posts .
ENTRYPOINT ["./posts"]
