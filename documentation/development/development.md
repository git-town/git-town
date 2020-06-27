# Developing the Git Town source code

## setup

You need to have these things running on your computer:

- [Go](https://golang.org) version 1.9 or higher
- [Make](https://www.gnu.org/software/make)
  - Mac and Linux users should be okay out of the box
  - Windows users should install
    [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

Optional:
- [Yarn](https://yarnpkg.com/)
- [scc](https://github.com/boyter/scc)

Fork Git Town and clone your fork into a directory outside your GOPATH. Git Town
uses Go modules and doesn't work properly inside the GOPATH.

Cd into the directory you just cloned and run
<code textrun="verify-make-command">make setup</code> to download the dependencies.

To make sure everything works:

- build the tool: <code textrun="verify-make-command">make build</code>
  - now you have `$GOPATH/bin/git-town` compiled from your local source code
- run the tests: <code textrun="verify-make-command">make test</code>

## add a new Go library

- start using the new dependency in the code
- run `go mod vendor` to vendor it

## update a dependency

- `go get <path>`

## update all dependencies

<code textrun="verify-make-command">make update</code>

## auto-fix linter errors

```bash
make fix
```

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

### debug

To see the CLI output of the shell commands in a Cucumber test,
add a tag `@debug` above the feature or scenario you want to debug:

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
