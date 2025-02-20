FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . .
RUN PATH="/go/bin:${PATH}" GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go mod tidy && \
    go mod download && \
    go build -tags musl -ldflags '-s -w -extldflags "-static"' -o app ./cmd/auction/main.go

FROM scratch
COPY --from=builder /build/app ./app
EXPOSE 8080
ENTRYPOINT ["/app"]