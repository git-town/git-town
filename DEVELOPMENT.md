# Developing the Git Town source code

This page helps you get started hacking on the Git Town codebase.

### setup

1. install [Go](https://golang.org) version 1.16 or higher
2. install [Make](https://www.gnu.org/software/make)
   - Mac and Linux users have this out of the box
   - Windows users can install
     [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm) or
     run `choco install make` if [Chocolatey](https://chocolatey.org) is
     available.
3. run <code textrun="verify-make-command">make setup</code> to download the
   tooling needed to hack on this codebase
4. to make sure everything works:
   - compile the tool: <code textrun="verify-make-command">make build</code>
   - run the tests: <code textrun="verify-make-command">make test</code>
5. run all tests: `make test`

Optional dependencies:

- [dprint](https://dprint.dev)
- [Node.js](https://nodejs.org)
- [Yarn](https://yarnpkg.com/)
- [scc](https://github.com/boyter/scc)

### dependencies

- run `go get [dependency]` inside the Git Town folder to register the
  dependency
- use the new dependency in the code
- run `go mod vendor` to vendor it
- run `go mod tidy` to clean up

Update a Go dependency:

```
go get <path>
```

Update all dependencies:

```bash
make update
```

### unit tests

Run all unit tests:

```bash
make unit
```

Run unit tests for packages containing changes:

```
make u
```

Run an individual unit test:

```
go test src/cmd/root_test.go
go test src/cmd/root_test.go -v -run TestIsAcceptableGitVersion
```

### end-to-end tests

Run individual Cucumber tests:

```bash
make cuke                       # runs all end-to-end test
godog [path to file/folder]     # runs the given end-to-end tests
```

Certain tests require that the Git remote points to an actual GitHub, Gitea,
GitLab or Bitbucket address. This causes `git push` operations in this test to
also go to GitHub. To prevent this, set an environment variable
`GIT_TOWN_REMOTE` with the desired value of the `origin` remote, and Git Town
will use that value instead of what is in the repo.

If Cucumber tests produce garbled output on Windows, try running them inside Git
Bash. See [this issue](https://github.com/cucumber/godog/issues/129) for
details.

### debug end-to-end tests

To see the CLI output of the shell commands in a Cucumber test, add a tag
`@debug` above the feature or scenario you want to debug:

```cucumber
@debug
Scenario: A foo walks into a bar
  Given ...
```

Debug a Godog Cucumber feature in [VSCode](https://code.visualstudio.com):

- open `main_test.go`
- change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

### run linters

```bash
make lint     # run all linters
make fix      # auto-fix linter errors
```

### learn about the code and test architecture

The source code contains
[comments](https://pkg.go.dev/github.com/git-town/git-town) that explain the
code architecture.
