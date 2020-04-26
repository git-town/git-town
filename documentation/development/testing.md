# Testing

Git Town has a very thorough test suite to ensure it works correctly and never
loses your data. The end-to-end tests use [Cucumber](https://cucumber.io) to
make them easy to read and reason about. Unit tests are written as normal Go
tests.

## Running Tests

```bash
# running the different test types
make test       # runs all tests
make test-go    # runs the new Go-based test suite
make lint       # runs the linters
make cuke       # runs the feature tests

# running individual scenarios/features
godog <path to file/folder>
```

## Auto-fixing issues

Auto-fix formatting and some linter errors by running:

```bash
make fix
```

## Debugging

To see the CLI output of Cucumber tests, add a tag `@debug` above the feature or
scenario you want to debug. Here is an example:

```cucumber
@debug
Feature: A foo walks into a bar
```

To debug a Godog Cucumber feature in [VSCode](https://code.visualstudio.com):

- open `main_test.go`
- change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

## Mocking

Certain tests require the Git remote to be set to a real value on GitHub or
Bitbucket. This causes `git push` operations in this test to also go to GitHub,
which is undesirable. To prevent this problem, Git Town mocks the Git remote if
a Git configuration value `git-town.testing.remote-url` exists with the
respective value. This also requires an environment variable `GIT_TOWN_ENV` set
to `test`.

## Architecture

See the [test architecture](test-architecture.md) document
