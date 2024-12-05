.PHONY: test gherkin bump cover

VERS ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MINIMUM_SUPPORTED_GO_MAJOR_VERSION = 1
MINIMUM_SUPPORTED_GO_MINOR_VERSION = 16
GO_VERSION_VALIDATION_ERR_MSG = Go version $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION) is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)

.PHONY: check-go-version
check-go-version:
	@if [ $(GO_MAJOR_VERSION) -gt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MINOR_VERSION) ] ; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	fi

test: check-go-version
	@echo "running all tests"
	@go fmt ./...
	@go run honnef.co/go/tools/cmd/staticcheck@v0.4.7 github.com/cucumber/godog
	@go run honnef.co/go/tools/cmd/staticcheck@v0.4.7 github.com/cucumber/godog/cmd/godog
	go vet ./...
	go test -race ./...
	go run ./cmd/godog -f progress -c 4

gherkin:
	@if [ -z "$(VERS)" ]; then echo "Provide gherkin version like: 'VERS=commit-hash'"; exit 1; fi
	@rm -rf gherkin
	@mkdir gherkin
	@curl -s -L https://github.com/cucumber/gherkin-go/tarball/$(VERS) | tar -C gherkin -zx --strip-components 1
	@rm -rf gherkin/{.travis.yml,.gitignore,*_test.go,gherkin-generate*,*.razor,*.jq,Makefile,CONTRIBUTING.md}

bump:
	@if [ -z "$(VERSION)" ]; then echo "Provide version like: 'VERSION=$(VERS) make bump'"; exit 1; fi
	@echo "bumping version from: $(VERS) to $(VERSION)"
	@sed -i.bak 's/$(VERS)/$(VERSION)/g' godog.go
	@sed -i.bak 's/$(VERS)/$(VERSION)/g' _examples/api/features/version.feature
	@find . -name '*.bak' | xargs rm

cover:
	go test -race -coverprofile=coverage.txt
	go tool cover -html=coverage.txt
	rm coverage.txt

ARTIFACT_DIR := _artifacts

# To upload artifacts for the current version;
# execute: make upload
#
# Check https://github.com/tcnksm/ghr for usage of ghr
upload: artifacts
	ghr -replace $(VERS) $(ARTIFACT_DIR)

# To build artifacts for the current version;
# execute: make artifacts
artifacts: 
	rm -rf $(ARTIFACT_DIR)
	mkdir $(ARTIFACT_DIR)

	$(call _build,darwin,amd64)
	$(call _build,linux,amd64)
	$(call _build,linux,arm64)

define _build
	mkdir $(ARTIFACT_DIR)/godog-$(VERS)-$1-$2
	env GOOS=$1 GOARCH=$2 go build -ldflags "-X github.com/cucumber/godog.Version=$(VERS)" -o $(ARTIFACT_DIR)/godog-$(VERS)-$1-$2/godog ./cmd/godog
	cp README.md $(ARTIFACT_DIR)/godog-$(VERS)-$1-$2/README.md
	cp LICENSE $(ARTIFACT_DIR)/godog-$(VERS)-$1-$2/LICENSE
	cd $(ARTIFACT_DIR) && tar -c --use-compress-program="pigz --fast" -f godog-$(VERS)-$1-$2.tar.gz godog-$(VERS)-$1-$2 && cd ..
	rm -rf $(ARTIFACT_DIR)/godog-$(VERS)-$1-$2
endef
