build:
    go build -o multiplayer-service ./cmd/main.go

run:
    go run ./cmd/main.go

test:
    go test ./...
