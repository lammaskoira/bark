## Settings

# Utility settings
TOOLS_DIR := .tools
GOLANGCI_LINT_VERSION = v1.46.2

## Targets

all: lint test
PHONY: test coverage lint golint clean vendor

.PHONY: test coverage lint golint vendor clean

build:
	go build -o bark main.go

test: | lint
	@echo Running unit tests...
	@go test -timeout 30s -cover -short  -tags testtools ./...

coverage:
	@echo Generating coverage report...
	@go test -timeout 30s -tags testtools ./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

golint: | vendor $(TOOLS_DIR)/golangci-lint
	@echo Linting Go files...
	@$(TOOLS_DIR)/golangci-lint run

clean:
	@echo Cleaning...
	@rm -rf coverage.out
	@go clean -testcache
	@rm -r $(TOOLS_DIR)

vendor:
	@go mod download
	@go mod tidy

# Tools setup
$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	export \
		VERSION=$(GOLANGCI_LINT_VERSION) \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=$(TOOLS_DIR) && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION
	$(TOOLS_DIR)/golangci-lint version
	$(TOOLS_DIR)/golangci-lint linters
