# Build Chat Bot app
FROM golang:1.13-buster as builder

WORKDIR /app

COPY . .

RUN go mod download && \
    go mod verify && \
    CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 \
        go build \
        -o app \
        main.go

# Copy chat bot app binary to executor image
FROM keybaseio/client:nightly-slim

COPY --from=builder /app/app /usr/bin/app
