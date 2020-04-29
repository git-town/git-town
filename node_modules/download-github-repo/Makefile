
node_modules: package.json
	@npm install

test: node_modules
	@node_modules/.bin/mocha --reporter spec --timeout 10000

.PHONY: test
