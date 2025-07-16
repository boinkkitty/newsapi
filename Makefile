fmt::
	go fmt ./...

run::
	go run ./cmd/api-server/main.go

tidy::
	go mod tidy -v

test::
	go test ./...