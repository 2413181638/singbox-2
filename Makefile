GOFLAGS ?=
BIN_DIR := bin
PACKAGE := github.com/yourusername/singxclient

.PHONY: build
build:
	go build $(GOFLAGS) -o $(BIN_DIR)/singxclient ./cmd/singxclient

.PHONY: test
test:
	go test ./...

.PHONY: generate
generate:
	go install github.com/goreleaser/goreleaser@latest