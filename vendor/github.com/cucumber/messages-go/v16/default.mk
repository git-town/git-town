# Please update /.templates/go/default.mk and sync:
#  source /scripts/functions.sh && rsync_files

SHELL := /usr/bin/env bash
GOPATH := $(shell go env GOPATH)
PATH := $(PATH):$(GOPATH)/bin
GO_SOURCE_FILES := $(shell find . -name "*.go" | sort)
LIBNAME := $(shell basename $$(dirname $$(pwd)))
EXE_BASE_NAME := cucumber-$(LIBNAME)
LDFLAGS := "-X main.version=${NEW_VERSION}"

# Enumerating Cross compilation targets
PLATFORMS = darwin-amd64 linux-386 linux-amd64 linux-arm freebsd-386 freebsd-amd64 openbsd-386 openbsd-amd64 windows-386 windows-amd64 freebsd-arm netbsd-386 netbsd-amd64 netbsd-arm
PLATFORM = $(patsubst dist/$(EXE_BASE_NAME)-%,%,$@)
OS_ARCH = $(subst -, ,$(PLATFORM))
X-OS = $(word 1, $(OS_ARCH))
X-ARCH = $(word 2, $(OS_ARCH))

# Determine if we're on linux or osx (ignoring other OSes as we're not building on them)
OS := $(shell [[ "$$(uname)" == "Darwin" ]] && echo "darwin" || echo "linux")
# Determine if we're on 386 or amd64 (ignoring other processors as we're not building on them)
ARCH := $(shell [[ "$$(uname -m)" == "x86_64" ]] && echo "amd64" || echo "386")
EXE := dist/$(EXE_BASE_NAME)-$(OS)-$(ARCH)

ifndef NO_CROSS_COMPILE
EXES = $(patsubst %,dist/$(EXE_BASE_NAME)-%,$(PLATFORMS))
else
EXES = $(EXE)
endif

GO_REPLACEMENTS := $(shell sed -n "/^\s*github.com\/cucumber/p" go.mod | perl -wpe 's/\s*(github.com\/cucumber\/(.*)-go\/v\d+).*/q{replace } . $$1 . q{ => ..\/..\/} . $$2 . q{\/go}/eg')
CURRENT_MAJOR := $(shell sed -n "/^module/p" go.mod | awk '{ print $$0 "/v1" }' | cut -d'/' -f4 | cut -d'v' -f2)
NEW_MAJOR := $(shell echo ${NEW_VERSION} | awk -F'.' '{print $$1}')

GO_MAJOR_V = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_V = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MIN_SUPPORTED_GO_MAJOR_V = 1
MIN_SUPPORTED_GO_MINOR_V = 13

# https://stackoverflow.com/questions/2483182/recursive-wildcards-in-gnu-make
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

default: .linted .tested
.PHONY: default

# Run the .dist target if there is a main file
ifneq (,$(wildcard ./cmd/main.go))
default: dist
endif

.deps:
	touch $@

dist: $(EXES)

dist/$(EXE_BASE_NAME)-%: .deps $(GO_SOURCE_FILES)
	mkdir -p dist
	echo "EXES=$(EXES)"
	echo "Building $@"

	# Determine if we're on a supported go platform
	@if [ $(GO_MAJOR_V) -gt $(MIN_SUPPORTED_GO_MAJOR_V) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_V) -lt $(MIN_SUPPORTED_GO_MAJOR_V) ]; then \
		echo '$(GO_MAJOR_V).$(GO_MINOR_V) is not a supported version, $(MIN_SUPPORTED_GO_MAJOR_V).$(MIN_SUPPORTED_GO_MINOR_V) is required';\
		exit 1; \
	elif [ $(GO_MINOR_V) -lt $(MIN_SUPPORTED_GO_MINOR_V) ] ; then \
		echo '$(GO_MAJOR_V).$(GO_MINOR_V) is not a supported version, $(MIN_SUPPORTED_GO_MAJOR_V).$(MIN_SUPPORTED_GO_MINOR_V) is required';\
		exit 1; \
	fi

	GOOS=$(X-OS) GOARCH=$(X-ARCH) go build -buildmode=exe -ldflags $(LDFLAGS) -o $@ -a ./cmd
ifndef NO_UPX_COMPRESSION
	# requires upx in PATH to compress supported binaries
	# may produce an error ARCH not supported
	-upx $@ -o $@.upx

	# If the compressed file passes the integrity test, replace the original
	# file with the compressed file. Otherwise, preserve the original file
	# and remove the compressed file.
	if [ -f "$@.upx" ]; then upx -t $@.upx && mv $@.upx $@ || rm -f $@.upx; fi
endif

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
	rm -rf .deps .tested* .linted dist/ acceptance/
.PHONY: clean-go

remove-replaces:
	sed -i '/^replace/d' go.mod
	sed -i 'N;/^\n$$/D;P;D;' go.mod
.PHONY: remove-replaces

add-replaces:
ifeq ($(shell sed -n "/^\s*github.com\/cucumber/p" go.mod | wc -l), 0)
	# No replacements here
else
	sed -i '/^go .*/i $(GO_REPLACEMENTS)\n' go.mod
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

### COMMON stuff for all platforms

BERP_VERSION = 1.3.0
BERP_GRAMMAR = gherkin.berp

define berp-generate-parser =
-! dotnet tool list --tool-path /usr/bin | grep "berp\s*$(BERP_VERSION)" && dotnet tool update Berp --version $(BERP_VERSION) --tool-path /usr/bin
berp -g $(BERP_GRAMMAR) -t $< -o $@ --noBOM
endef
