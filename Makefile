GO_BIN := $(shell pwd)/.bin

fmt::
	@export PATH=$(GO_BIN):$$PATH && golangci-lint run --fix -v ./...

run::
	go run ./cmd/api-server/main.go

tidy::
	go mod tidy -v

test::
	go test ./...

tools::
	mkdir -p $(GO_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_BIN) v1.64.5
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % sh -c 'GOBIN=${GO_BIN} go install %'