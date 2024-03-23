.PHONY: test gherkin bump cover

VERS ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

FOUND_GO_VERSION := $(shell go version)
EXPECTED_GO_VERSION = 1.17
.PHONY: check-go-version
check-go-version:
	@$(if $(findstring ${EXPECTED_GO_VERSION}, ${FOUND_GO_VERSION}),(exit 0),(echo Wrong go version! Please install ${EXPECTED_GO_VERSION}; exit 1))

test: check-go-version
	@echo "running all tests"
	@go install ./...
	@go fmt ./...
	@go run honnef.co/go/tools/cmd/staticcheck@v0.2.2 github.com/cucumber/godog
	@go run honnef.co/go/tools/cmd/staticcheck@v0.2.2 github.com/cucumber/godog/cmd/godog
	go vet ./...
	go test -race ./...
	godog -f progress -c 4

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
