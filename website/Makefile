.DEFAULT_GOAL := help

build:  # transpiles the website to HTML
	mdbook build

clean:  # removes all build artifacts
	mdbook clean

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

serve:  # runs a local development server of the website
	mdbook serve --open
