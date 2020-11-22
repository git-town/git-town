# Developing the Git Town source code

This page gets you started hacking on the Git Town codebase.

## setup

1. install [Go](https://golang.org) version 1.12 or higher
2. install [Make](https://www.gnu.org/software/make)
   - Mac and Linux users should be okay out of the box
   - Windows users can install
     [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm). If
     you use [Chocolatey](https://chocolatey.org), run `choco install make`.
3. create a fork of the
   [Git Town repository on GitHub](https://github.com/git-town/git-town) by
   clicking on `Fork` there
4. clone your fork into a directory outside your GOPATH. Git Town uses Go
   modules and doesn't work properly inside the GOPATH. If you don't know what a
   GOPATH is, clone into a directory other than `c:\go` and `~/go`.
5. open a terminal and cd into the directory you cloned
6. run <code textrun="verify-make-command">make setup</code> to download the
   dependencies
7. make sure everything works:
   - build the tool: <code textrun="verify-make-command">make build</code>
   - run the tests: <code textrun="verify-make-command">make test</code>

Optional dependencies:

- [Node.JS](https://nodejs.org)
- [Yarn](https://yarnpkg.com/)
- [scc](https://github.com/boyter/scc)

## add a new Go dependency

- run `go get [dependency]` inside the Git Town folder to register the
  dependency
- start using the new dependency in the code
- run `go mod vendor` to vendor it
- run `go mod tidy` to clean up

## update a dependency

```
go get <path>
```

## update all dependencies

<code textrun="verify-make-command">make update</code>

## auto-fix linter errors

<pre textrun="verify-make-command">
make fix
</pre>

## run tests

```bash
make test       # runs all tests
make test-go    # runs the Go tests (faster during development)
make cuke       # runs the feature tests
make lint       # runs the linters
```

Run individual Cucumber tests:

```bash
godog [path to file/folder]
```

Run individual unit tests:

```
go test src/cmd/root_test.go
go test src/cmd/root_test.go -v -run TestIsAcceptableGitVersion
```

Certain tests require that the Git remote points to an actual GitHub, Gitea,
GitLab or Bitbucket address. This causes `git push` operations in this test to
also go to GitHub. To prevent this, set an environment variable
`GIT_TOWN_REMOTE` with the desired value of the `origin` remote, and Git Town
will use that value instead of what is in the repo.

See the [test architecture](test-architecture.md) document for more details.

If Cucumber tests produce garbled output on Windows, please try running them
inside Git Bash. See [this issue](https://github.com/cucumber/godog/issues/129)
for details. The test suite doesn't run browser tests because the Windows
`start` CLI command is a built-in shell command and cannot be mocked. Tests
asking for user input are disabled as well because of problems piping input into
subshells on Windows. The business logic for these features are covered on
non-Windows machines.

### debug

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
