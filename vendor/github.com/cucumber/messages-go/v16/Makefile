include default.mk

JSONSCHEMAS = $(shell find ../jsonschema -name "*.json")

.deps: messages.go

messages.go: $(JSONSCHEMAS) ../jsonschema/scripts/codegen.rb
	ruby ../jsonschema/scripts/codegen.rb Go ../jsonschema > $@
