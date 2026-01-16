RTA_VERSION = 0.28.0  # run-that-app version to use

.DEFAULT_GOAL := help

build: ../tools/rta@${RTA_VERSION}  # transpiles the website to HTML
	../tools/rta --include=mdbook-linkcheck mdbook build

clean: ../tools/rta@${RTA_VERSION}  # removes all build artifacts
	../tools/rta mdbook clean
	rm -rf bin

help:  # prints available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

serve: ../tools/rta@${RTA_VERSION}  # runs a local development server of the website
	PATH=$(MDBOOK_LINKCHECK_PATH):$$PATH ../tools/rta mdbook serve --open -n 127.0.0.1

test: build unit  # tests the website
	cd .. && make --no-print-directory docs

unit: ../tools/rta@${RTA_VERSION}  # runs the unit tests
	cd ../tools && ./rta node --test

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

../tools/rta@${RTA_VERSION}:
	@rm -f ../tools/rta*
	@(cd ../tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh -s ${RTA_VERSION})
	@mv ../tools/rta ../tools/rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} ../tools/rta
