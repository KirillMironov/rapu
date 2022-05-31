FROM golang:1.17.5-alpine3.14 AS modules
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY users users

FROM modules AS tester
WORKDIR /app/users
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go test -v ./...

FROM modules AS builder
WORKDIR /app/users
ARG TARGETARCH
ARG TARGETOS
RUN GOARCH=$TARGETARCH GOOS=$TARGETOS CGO_ENABLED=0 go build .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/users/users .
ENTRYPOINT ["./users"]
