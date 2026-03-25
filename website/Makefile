RTA_VERSION = 0.30.0  # run-that-app version to use

.DEFAULT_GOAL := help

RTA     = ../tools/rta@$(RTA_VERSION)
MDBOOK  = $(RTA) mdbook

build: ${RTA}  # transpiles the website to HTML
	$(RTA) --include=mdbook-linkcheck mdbook build

clean: ${RTA}  # removes all build artifacts
	$(MDBOOK) clean
	rm -rf bin

help:  # prints available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

serve: ${RTA}  # runs a local development server of the website
	PATH=$(MDBOOK_LINKCHECK_PATH):$$PATH $(MDBOOK) serve --open -n 127.0.0.1

test: build unit  # tests the website
	cd .. && make --no-print-directory docs

unit: ${RTA}  # runs the unit tests
	cd ../tools && ./rta@${RTA_VERSION} node --test

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

../tools/rta@${RTA_VERSION}:
	@rm -f ../tools/rta*
	@(cd ../tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh -s -- --version ${RTA_VERSION} --name rta@${RTA_VERSION})
