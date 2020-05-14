# Testing

Git Town has a comprehensive test suite to ensure it never loses data. The
end-to-end tests use [Cucumber](https://cucumber.io) to make them easier to
reason about. Unit tests are normal Go tests.

## Running Tests

To run different test types:

```bash
make test       # runs all tests
make test-go    # runs the new Go-based test suite
make lint       # runs the linters
make cuke       # runs the feature tests
```

To run individual scenarios/features:

```bash
godog [path to file/folder]
```

## Auto-fixing issues

Auto-fix formatting and some linter errors by running:

```bash
make fix
```

## Debugging

To see the CLI output of commands in Cucumber tests, add a tag `@debug` above
the feature or scenario you want to debug. Here is an example:

```cucumber
@debug
Scenario: A foo walks into a bar
  Given ...
```

To debug a Godog Cucumber feature in [VSCode](https://code.visualstudio.com):

- open `main_test.go`
- change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

## Preventing pushes to GitHub

Certain tests require that the Git remote points to an actual GitHub or
Bitbucket address. This causes `git push` operations in this test to also go to
GitHub. To prevent this, set an environment variable `GIT_TOWN_REMOTE` with the
desired value of the `origin` remote, and Git Town will use that value instead
of the Git Town configuration.

## Architecture

See the [test architecture](test-architecture.md) document
