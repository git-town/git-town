##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

reviewable: setup generate fmt lint test ## Run before committing.

fmt: ## Format code
	@buf format -w
	@gofumpt -l -w *.go testing/*.go examples/*.go

lint: ## Run linter
	@golangci-lint run
	@buf format --exit-code
	@buf lint

.PHONY: setup
setup: ## Setup your local environment
	go mod tidy

.PHONY: generate
generate: ## Generate files
	buf generate # install from .tool-versions
	./scripts/generate_testing_client.sh
	./scripts/generate_service_interface_map.sh
	./scripts/generate_mock_api.sh

.PHONY: clean
clean: ## Remove generated files
	rm -f \
		testing/*_mock.go \
		testing/*_generated.go \
		*_generated_test.go

test: ## Run tests
	go test ./... -race
