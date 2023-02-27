# the version of mdbook to use
MDBOOK_VERSION = 0.4.15


# the full path to the mdbook binary to use
mdbook_path = bin/mdbook_${MDBOOK_VERSION}


build: ${mdbook_path}  # transpiles the website to HTML
	${mdbook_path} build

clean: ${mdbook_path}  # removes all build artifacts
	${mdbook_path} clean
	rm -rf bin

help:  # prints available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v help | grep -v '^bin\/' | grep -v '$${mdbook_path}:' | grep -v '.SILENT' | sed 's/:.*#/#/' | column -s "#" -t

serve: ${mdbook_path}  # runs a local development server of the website
	${mdbook_path} serve --open

test:  # tests the website
	cd .. && make --no-print-dir docs

${mdbook_path}:  # installs the mdbook binary
	scripts/install_mdbook ${MDBOOK_VERSION}


.DEFAULT_GOAL := help
.SILENT:
