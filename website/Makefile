build: bin/mdbook  # transpiles the website to HTML
	bin/mdbook build

clean: bin/mdbook  # removes all build artifacts
	bin/mdbook clean
	rm -rf bin

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | grep -v '^bin\/' | grep -v '.SILENT' | sed 's/:.*#/#/' | column -s "#" -t

serve: bin/mdbook  # runs a local development server of the website
	bin/mdbook serve --open

test:  # tests the website
	cd .. && make --no-print-dir docs

bin/mdbook: scripts/install_mdbook  # installs the mdBook binary
	scripts/install_mdbook


.DEFAULT_GOAL := help
.SILENT:
