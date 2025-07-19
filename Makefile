GO_BIN := $(shell pwd)/.bin

fmt::
	@export PATH=$(GO_BIN):$$PATH && golangci-lint run --fix -v ./...

generate::
	@PATH=$(GO_BIN):$$PATH go generate ./...

run::
	go run ./cmd/api-server/main.go

tidy::
	go mod tidy -v

test::
	go test ./...

tools::
	mkdir -p $(GO_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) v1.64.5
	GOBIN=${GO_BIN} go install tool