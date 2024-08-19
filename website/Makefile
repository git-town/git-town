RTA_VERSION = 0.6.1  # run-that-app version to use

build: ../tools/rta@${RTA_VERSION}  # transpiles the website to HTML
	../tools/rta mdbook build

clean: ../tools/rta@${RTA_VERSION} # removes all build artifacts
	../tools/rta mdbook clean
	rm -rf bin

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | grep -v '^bin\/' | grep -v '$${mdbook_path}:' | grep -v '.SILENT' | sed 's/:.*#/#/' | column -s "#" -t

serve: ../tools/rta@${RTA_VERSION}  # runs a local development server of the website
	../tools/rta mdbook serve --open

test:  # tests the website
	cd .. && make --no-print-dir docs

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

../tools/rta@${RTA_VERSION}:
	@rm -f ../tools/rta*
	@(cd ../tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv ../tools/rta ../tools/rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} ../tools/rta

.DEFAULT_GOAL := help
