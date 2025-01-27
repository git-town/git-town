RTA_VERSION = 0.10.4  # run-that-app version to use

.DEFAULT_GOAL := help
MDBOOK_LINKCHECK_PATH := $(shell dirname $(shell rta --which mdbook-linkcheck))

build: rta@${RTA_VERSION}  # transpiles the website to HTML
	PATH=$(MDBOOK_LINKCHECK_PATH):$$PATH rta mdbook build

clean: rta@${RTA_VERSION}  # removes all build artifacts
	rta mdbook clean
	rm -rf bin

help:  # prints available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

serve: rta@${RTA_VERSION}  # runs a local development server of the website
	PATH=$(MDBOOK_LINKCHECK_PATH):$$PATH rta mdbook serve --open -n 127.0.0.1

test: build unit  # tests the website
	cd .. && make --no-print-directory docs

unit: rta@${RTA_VERSION}  # runs the unit tests
	cd rta node --test

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

rta@${RTA_VERSION}:
	@rm -f rta*
	@(curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv rta rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} rta
