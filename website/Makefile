.DEFAULT_GOAL := help

build:  # transpiles the website to HTML
	./mdbook build

clean:  # removes all build artifacts
	./mdbook clean
	rm -rf bin

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

serve:  # runs a local development server of the website
	./mdbook serve --open

setup:  # installs the mdBook binary
	./scripts/install_mdbook

test:  # tests the website
	cd .. && make --no-print-dir docs


.SILENT:
