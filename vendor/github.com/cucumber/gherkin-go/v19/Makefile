include default.mk

GHERKIN_DIALECTS := $(shell cat gherkin-languages.json | jq --compact-output --sort-keys . | base64 | tr -d '\n')
LDFLAGS := "-X 'main.version=${NEW_VERSION}' -X 'main.gherkinDialects=${GHERKIN_DIALECTS}'"

GOOD_FEATURE_FILES = $(shell find testdata/good -name "*.feature")
BAD_FEATURE_FILES  = $(shell find testdata/bad -name "*.feature")

TOKENS_GOLDEN       = $(patsubst testdata/%.feature,testdata/%.feature.tokens,$(GOOD_FEATURE_FILES))
ASTS_GOLDEN         = $(patsubst testdata/%.feature,testdata/%.feature.ast.ndjson,$(GOOD_FEATURE_FILES))
PICKLES_GOLDEN      = $(patsubst testdata/%.feature,testdata/%.feature.pickles.ndjson,$(GOOD_FEATURE_FILES))
SOURCES_GOLDEN      = $(patsubst testdata/%.feature,testdata/%.feature.source.ndjson,$(GOOD_FEATURE_FILES))
ERRORS_GOLDEN       = $(patsubst testdata/%.feature,testdata/%.feature.errors.ndjson,$(BAD_FEATURE_FILES))

TOKENS       = $(patsubst testdata/%.feature,acceptance/testdata/%.feature.tokens,$(GOOD_FEATURE_FILES))
ASTS         = $(patsubst testdata/%.feature,acceptance/testdata/%.feature.ast.ndjson,$(GOOD_FEATURE_FILES))
PICKLES      = $(patsubst testdata/%.feature,acceptance/testdata/%.feature.pickles.ndjson,$(GOOD_FEATURE_FILES))
SOURCES      = $(patsubst testdata/%.feature,acceptance/testdata/%.feature.source.ndjson,$(GOOD_FEATURE_FILES))
ERRORS       = $(patsubst testdata/%.feature,acceptance/testdata/%.feature.errors.ndjson,$(BAD_FEATURE_FILES))

.DELETE_ON_ERROR:

default: .compared

.compared: bin/gherkin-generate-tokens $(EXE) $(TOKENS) $(ASTS) $(PICKLES) $(SOURCES) $(ERRORS)
	touch $@

.tested: parser.go dialects_builtin.go

# The golden target regenerates the golden master
golden: bin/gherkin-generate-tokens $(EXE) $(TOKENS_GOLDEN) $(ASTS_GOLDEN) $(PICKLES_GOLDEN) $(SOURCES_GOLDEN) $(ERRORS_GOLDEN)

bin/gherkin-generate-tokens: .deps $(GO_SOURCE_FILES) parser.go dialects_builtin.go
	go build -o $@ ./gherkin-generate-tokens

$(EXE): parser.go dialects_builtin.go

testdata/%.feature.tokens: testdata/%.feature
ifdef GOLDEN
	mkdir -p $(@D)
	bin/gherkin-generate-tokens $< > $@
endif

acceptance/testdata/%.feature.tokens: testdata/%.feature testdata/%.feature.tokens bin/gherkin-generate-tokens
	mkdir -p $(@D)
	bin/gherkin-generate-tokens $< > $@
	diff --unified $<.tokens $@

testdata/%.feature.ast.ndjson: testdata/%.feature $(EXE)
ifdef GOLDEN
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source --no-pickles $< | jq --sort-keys  --compact-output > $@
endif

acceptance/testdata/%.feature.ast.ndjson: testdata/%.feature testdata/%.feature.ast.ndjson $(EXE)
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source --no-pickles $< | jq --sort-keys --compact-output > $@
	diff --unified <(jq "." $<.ast.ndjson) <(jq "." $@)

testdata/%.feature.errors.ndjson: testdata/%.feature $(EXE)
ifdef GOLDEN
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source $< | jq --sort-keys  --compact-output > $@
endif

acceptance/testdata/%.feature.errors.ndjson: testdata/%.feature testdata/%.feature.errors.ndjson $(EXE)
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source $< | jq --sort-keys --compact-output > $@
	diff --unified <(jq "." $<.errors.ndjson) <(jq "." $@)

testdata/%.feature.source.ndjson: testdata/%.feature $(EXE)
ifdef GOLDEN
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-ast --no-pickles $< | jq --sort-keys  --compact-output > $@
endif

acceptance/testdata/%.feature.source.ndjson: testdata/%.feature testdata/%.feature.source.ndjson $(EXE)
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-ast --no-pickles $< | jq --sort-keys --compact-output > $@
	diff --unified <(jq "." $<.source.ndjson) <(jq "." $@)

testdata/%.feature.pickles.ndjson: testdata/%.feature $(EXE)
ifdef GOLDEN
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source --no-ast $< | jq --sort-keys --compact-output > $@
endif

acceptance/testdata/%.feature.pickles.ndjson: testdata/%.feature testdata/%.feature.pickles.ndjson $(EXE)
	mkdir -p $(@D)
	$(EXE) --predictable-ids --no-source --no-ast $< | jq --sort-keys --compact-output > $@
	diff --unified <(jq "." $<.pickles.ndjson) <(jq "." $@)

parser.go: parser.go.razor gherkin.berp
	$(berp-generate-parser)
	gofmt -w $@

dialects_builtin.go: gherkin-languages.json dialects_builtin.go.jq
	cat $< | jq --sort-keys --from-file dialects_builtin.go.jq --raw-output --compact-output > $@
	gofmt -w $@

clean:
	rm -rf .compared bin/
