# Developing the Git Town source code

This page helps you get started hacking on the Git Town codebase.

## setup

1. install [Go](https://golang.org) version 1.16 or higher
2. install [Make](https://www.gnu.org/software/make)
   - Mac and Linux users have this out of the box
   - Windows users can install
     [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm) or
     run `choco install make` if [Chocolatey](https://chocolatey.org) is
     available.
3. automatically install Go-based tooling
   <code type="make/command">make setup</code>
4. manually install optional tooling: [dprint](https://dprint.dev),
   [Node.js](https://nodejs.org), [Yarn](https://yarnpkg.com/),
   [scc](https://github.com/boyter/scc)
5. run the tests: <code type="make/command">make test</code>
   * Note: When running on a non-English locale, run with `LANG=C make test` to prevent tests from failing.
6. compile the tool: <code type="make/command">make build</code>

## dependencies

Add an external Go module:

- run `go get [dependency]` inside the Git Town folder to register the
  dependency
- use the new dependency in the code
- run `go mod vendor` to vendor it
- run `go mod tidy` to clean up

Update an external Go module:

```
go get <path>
```

Update all external Go modules:

<a type="make/command">

```
make update
```

</a>

## unit tests

Run all unit tests:

<a type="make/command">

```
make unit
```

</a>

Run unit tests for packages containing changes:

<a type="make/command">

```
make u
```

</a>

Run an individual unit test:

```
go test src/cmd/root_test.go
go test src/cmd/root_test.go -v -run TestIsAcceptableGitVersion
```

## end-to-end tests

Run all end-to-end tests:

<a type="make/command">

```
make cuke
```

</a>

Run individual Cucumber tests:

```bash
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

## debug end-to-end tests

To see the CLI output of the shell commands in a Cucumber test, add a tag
`@debug` above the feature or scenario you want to debug:

```cucumber
@debug
Scenario: A foo walks into a bar
```

To inspect the workspace containing the various Git repositories at any step of
an end-to-end test, insert the step `And inspect the repo`. Cucumber will print
the path of the workspace and wait until you hit ENTER.

Debug a Godog Cucumber feature in [VSCode](https://code.visualstudio.com):

- open `main_test.go`
- change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

## run linters

Run all linters:

<a type="make/command">

```
make lint
```

</a>

Auto-fix linter errors:

<a type="make/command">

```
make fix
```

</a>

## learn about the code and test architecture

The source code contains
[Godoc comments](https://pkg.go.dev/github.com/git-town/git-town) that explain
the code architecture.

## website

The source code for the [website](https://www.git-town.com) is in the
[website](website) folder. This folder contains its own
[Makefile](website/Makefile) for activities related to working on the website.
To work on the website, cd into the `website` folder and run
<code type="make/command" dir="website">make setup</code> to download the
necessary tooling and then
<code type="make/command" dir="website">make serve</code> to start a local
development server. The production site auto-updates on changes to the `main`
branch. The site hoster is [Netlify](https://www.netlify.com). Netlify
configuration is in [netlify.toml](netlify.toml).
