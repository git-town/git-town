# AssertNoDiff

> Asserts the equality of JavaScript strings and objects with human-friendly,
> Bash-colored diffs as error messages.

[![CircleCI](https://circleci.com/gh/kevgo/assert-no-diff.svg?style=shield)](https://circleci.com/gh/kevgo/assert-no-diff)
[![Coverage Status](https://coveralls.io/repos/github/kevgo/assert-no-diff/badge.svg?branch=master)](https://coveralls.io/github/kevgo/assert-no-diff?branch=master)
[![Language grade: JavaScript](https://img.shields.io/lgtm/grade/javascript/g/kevgo/assert-no-diff.svg)](https://lgtm.com/projects/g/kevgo/assert-no-diff/context:javascript)

This library provides the JavaScript string and object diff methods of the
popular [diff](https://github.com/kpdecker/jsdiff) library as assertions. The
exceptions thrown by the assertions highlight the differences of the actual and
expected values with human-friendly Bash-colored messages. This is useful for
verifying the equality of larger strings or data structures in tests.

```javascript
const assertNoDiff = require("assert-no-diff")

// assert no differences between two JavaScript objects
assertNoDiff.json(actualJson, expectedJson)

// compare strings and highlight the differences character-by-character
assertNoDiff.chars(actualString, expectedString)

// compare strings and highlight the mismatching words, whitespace-sensitive
assertNoDiff.wordsWithSpace(actualString, expectedString)

// compare strings and highlight the mismatching lines, ignoring whitespace around them
assertNoDiff.trimmedLines(actualString, expectedString)
```

You can provide a custom error message as an optional third parameter.

## Development

- run all tests: `make test`
- see all available make commands: `make help`
- deploy a new version:
  - update `package.json` to the new version and commit to master
  - run `npm publish`
