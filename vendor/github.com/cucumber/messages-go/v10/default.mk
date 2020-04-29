SHELL := /usr/bin/env bash
GOPATH := $(shell go env GOPATH)
PATH := $(PATH):$(GOPATH)/bin
GO_SOURCE_FILES := $(shell find . -name "*.go" | sort)
LIBNAME := $(shell basename $$(dirname $$(pwd)))
EXE_BASE_NAME := cucumber-$(LIBNAME)
GOX_LDFLAGS := "-X main.version=${NEW_VERSION}"
EXES := $(shell find dist -name '$(EXE_BASE_NAME)-*')
UPX_EXES = $(patsubst dist/$(EXE_BASE_NAME)-%,dist_compressed/$(EXE_BASE_NAME)-%,$(EXES))
# Determine if we're on linux or osx (ignoring other OSes as we're not building on them)
OS := $(shell [[ "$$(uname)" == "Darwin" ]] && echo "darwin" || echo "linux")
# Determine if we're on 386 or amd64 (ignoring other processors as we're not building on them)
ARCH := $(shell [[ "$$(uname -m)" == "x86_64" ]] && echo "amd64" || echo "386")
EXE = dist/$(EXE_BASE_NAME)-$(OS)-$(ARCH)
REPLACEMENTS := $(shell sed -n "/^\s*github.com\/cucumber/p" go.mod | perl -wpe 's/\s*(github.com\/cucumber\/(.*)-go\/v\d+).*/q{replace } . $$1 . q{ => ..\/..\/} . $$2 . q{\/go}/eg')
CURRENT_MAJOR := $(shell sed -n "/^module/p" go.mod | awk '{ print $$0 "/v1" }' | cut -d'/' -f4 | cut -d'v' -f2)
NEW_MAJOR := $(shell echo ${NEW_VERSION} | awk -F'.' '{print $$1}')

default: .linted .tested
.PHONY: default

# Run the .dist target if there is a main file
ifneq (,$(wildcard ./cmd/main.go))
default: dist
endif

.deps:
	touch $@

dist: $(EXE)
ifndef NO_UPX_COMPRESSION
	make .dist-compressed
endif
	touch $@

$(EXE): .deps $(GO_SOURCE_FILES)
	mkdir -p dist
ifndef NO_CROSS_COMPILE
	# Cross-compile executable for many platforms
	go get github.com/aslakhellesoy/gox
	gox -buildmode=exe -ldflags $(GOX_LDFLAGS) -output "dist/$(EXE_BASE_NAME)-{{.OS}}-{{.Arch}}" -rebuild ./cmd
else
	# Compile executable for the local platform only
	go build -ldflags $(GOX_LDFLAGS) -o $@ ./cmd
endif

.dist-compressed: $(UPX_EXES)
	touch $@

update-dependencies:
	go get -u && go mod tidy
.PHONY: update-dependencies

pre-release: remove-replaces update-version update-dependencies clean default
.PHONY: pre-release

update-version: update-major
	# no-op
.PHONY: update-version

ifneq (,$(wildcard ./cmd/main.go))
publish: dist
ifdef NEW_VERSION
	./scripts/github-release $(NEW_VERSION)
else
	@echo -e "\033[0;31mNEW_VERSION is not defined. Can't publish :-(\033[0m"
	exit 1
endif
else
publish:
	# no-op
endif
.PHONY: publish

dist_compressed/$(EXE_BASE_NAME)-%: dist/$(EXE_BASE_NAME)-%
	mkdir -p dist_compressed
	# requires upx in PATH to compress supported binaries
	# may produce an error ARCH not supported
	-upx $< -o $@

	# Test the integrity
	if [ -f "$@" ]; then upx -t $@ && cp $@ $< || rm $@; fi

.linted: $(GO_SOURCE_FILES)
	gofmt -w $^
	touch $@

.tested: .deps $(GO_SOURCE_FILES)
	go test ./...
	touch $@

post-release: add-replaces
.PHONY: post-release

clean: clean-go
.PHONY: clean

clean-go:
	rm -rf .deps .tested* .linted dist/ .dist-compressed dist_compressed/ acceptance/
.PHONY: clean-go

remove-replaces:
	sed -i '/^replace/d' go.mod
	sed -i 'N;/^\n$$/D;P;D;' go.mod
.PHONY: remove-replaces

add-replaces:
ifeq ($(shell sed -n "/^\s*github.com\/cucumber/p" go.mod | wc -l), 0)
	# No replacements here
else
	sed -i '/^go .*/i $(REPLACEMENTS)\n' go.mod
endif
.PHONY: add-replaces

update-major:
ifeq ($(CURRENT_MAJOR), $(NEW_MAJOR))
	# echo "No major version change"
else
	echo "Updating major from $(CURRENT_MAJOR) to $(NEW_MAJOR)"
	sed -Ei "s/$(LIBNAME)-go(\/v$(CURRENT_MAJOR))?/$(LIBNAME)-go\/v$(NEW_MAJOR)/" go.mod
	sed -Ei "s/$(LIBNAME)-go(\/v$(CURRENT_MAJOR))?/$(LIBNAME)-go\/v$(NEW_MAJOR)/" $(shell find . -name "*.go")
endif
.PHONY: update-major
