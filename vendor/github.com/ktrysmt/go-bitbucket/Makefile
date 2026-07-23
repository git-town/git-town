.DEFAULT_GOAL := help

env: ## check env for e2e testing
ifndef BITBUCKET_TEST_USERNAME
	$(error `BITBUCKET_TEST_USERNAME` is not set)
endif
ifndef BITBUCKET_TEST_PASSWORD
	$(error `BITBUCKET_TEST_PASSWORD` is not set)
endif
ifndef BITBUCKET_TEST_OWNER
	$(error `BITBUCKET_TEST_OWNER` is not set)
endif
ifndef BITBUCKET_TEST_REPOSLUG
	$(error `BITBUCKET_TEST_REPOSLUG` is not set)
endif
ifndef BITBUCKET_TEST_ACCESS_TOKEN
	$(error `BITBUCKET_TEST_ACCESS_TOKEN` is not set)
endif

build: ## compile all packages
	go build ./...

test: env test/unit test/mock test/e2e ## run all tests (requires env vars)

test/ci: build test/unit-short test/mock ## run tests that do not require credentials (for CI)

test/unit: ## run unit tests (httptest-based, no credentials needed)
	go test -v -count=1 .

test/unit-short: ## run unit tests without network-dependent tests
	go test -v -short -count=1 .

test/e2e: env ## run integration tests (requires Bitbucket credentials)
	go test -v ./tests

test/swagger: ## run integration tests against swagger mock server
	env BITBUCKET_API_BASE_URL=http://0.0.0.0:4010 go test -v ./tests

test/mock: ## run interface mock tests
	go test -v ./mock_tests

test/coverage: ## run unit tests with coverage report
	go test -coverprofile=coverage.out .
	go tool cover -func=coverage.out
	@echo "To view HTML report: go tool cover -html=coverage.out"

help: ## print this help
	@grep -E '^[a-zA-Z_/]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build test test/ci test/unit test/unit-short test/e2e test/swagger test/mock test/coverage help
