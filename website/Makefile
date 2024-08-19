# the version of mdbook to use
MDBOOK_VERSION = 0.4.15

build:  # transpiles the website to HTML
	../tools/rta mdbook build

clean:  # removes all build artifacts
	../tools/rta mdbook clean
	rm -rf bin

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | grep -v '^bin\/' | grep -v '$${mdbook_path}:' | grep -v '.SILENT' | sed 's/:.*#/#/' | column -s "#" -t

serve: ${mdbook_path}  # runs a local development server of the website
	${mdbook_path} serve --open

test:  # tests the website
	cd .. && make --no-print-dir docs


.DEFAULT_GOAL := help
.SILENT:
