RTA_VERSION = 0.6.1  # run-that-app version to use

fix: ../rta@${RTA_VERSION}  # runs all linters and auto-fixes
	tools/rta gofumpt -l -w .


help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: unit  # runs all tests
	../rta golangci-lint run
	../rta deadcode github.com/git-town/git-town/tools/structs_sorted

unit:  # runs only the unit tests
	go test ./...

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

tools/rta@${RTA_VERSION}:
	@rm -f ../rta*
	@(cd .. && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv ../rta ../rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} ../rta
