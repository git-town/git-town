# Developing the Git Town source code

This page helps you get started hacking on the Git Town codebase. See file
[ARCHITECTURE.md](ARCHITECTURE.md) for an overview of how the Git Town engine
works.

## setup

1. install [Go](https://golang.org) version 1.21
2. install [Make](https://www.gnu.org/software/make)
   - Mac and Linux users have this out of the box
   - Windows users can install
     [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm) or
     run `choco install make` if [Chocolatey](https://chocolatey.org) is
     available.
3. run all tests: <code type="make/command">make test</code>
4. install the tool locally: <code type="make/command">make build</code>
5. run a quick test suite during development: `make test-go`

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
make unit-all
```

</a>

Run all unit tests with race detection:

<a type="make/command">

```
make unit-race
```

</a>

Run unit tests for packages containing changes:

<a type="make/command">

```
make unit
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

To run individual Cucumber tests, add a `@this` flag to the test you want to
run. Example:

```cucumber
@this
Scenario: foo bar
```

Then run:

```
make cukethis
```

Certain tests require that the Git remote points to an actual GitHub, Gitea,
GitLab or Bitbucket address. This causes `git push` operations in this test to
also go to GitHub. To prevent this, set an environment variable
`GIT_TOWN_REMOTE` with the desired value of the `origin` remote, and Git Town
will use that value instead of what is in the repo.

If Cucumber tests produce garbled output on Windows, try running them inside Git
Bash. See [this issue](https://github.com/cucumber/godog/issues/129) for
details.

## inspecting variables

Inspect basic variables:

```go
fmt.Printf("%#v\n", variable)
```

Inspect more complex variables:

```go
import "github.com/davecgh/go-spew/spew"

spew.Dump(variable)
```

## debug end-to-end tests

To see the CLI output of the shell commands in a Cucumber test, as well as the
Git commands that the Git Town test suite runs under the hood, add a tag
`@debug` above the feature or scenario you want to debug:

```cucumber
@debug @this
Scenario: A foo walks into a bar
```

To see all Git commands that the test runner and the Git Town command execute,
runs the Git Town command with the `--verbose` option. As an example, if the
step `When I run "git-town append new"` mysteriously fails, you could change it
to `When I run "git-town append new -v"`. Also enable `@debug` to see the output
on the console.

To get a quick glance of which status the repo is at any point in time, insert
the step `And display "<command_>"` running whatever command you want to execute
in the Git repo under test. Example: `And display "git status"`.

To manually inspect the local and remote Git repositories used in an end-to-end
test, insert the step `And inspect the repo`. This step will make Cucumber print
the path of the workspace and wait until you hit ENTER.

Debug a Godog Cucumber feature in [VSCode](https://code.visualstudio.com):

- open `main_test.go`
- change the path of the test to execute
- set a breakpoint in your test code
- run the `debug a test` configuration in the debugger

## triangulate a hanging end-to-end test

End-to-end tests sometimes hang due to Git Town waiting for input that the test
doesn't define. To find the hanging test:

- open `main_test.go`
- find the call to `godog.RunWithOptions` and adjust its arguments:
  - Make the `Paths` field more specific, for example by changing it from
    "features" to "features/sync/current_branch". Now it runs only the tests in
    that subfolder.
  - To see the executed steps in the output, change `Format` to `pretty` and
    `Concurrency` to 1. This reduces the speed at which the end-to-end tests
    execute.

## run linters

Format all code, auto-fix all fixable issues, and run all linters:

<a type="make/command">

```
make fix
```

</a>

## debug the dialogs

Run `git town debug` to see the commands to display Git Town's dialogs.

## learn about the code and test architecture

The source code contains
[Godoc comments](https://pkg.go.dev/github.com/git-town/git-town) that explain
the code architecture.

## website

The source code for the [website](https://www.git-town.com) is in the
[website](website) folder. This folder contains its own
[Makefile](website/Makefile) for activities related to working on the website.
To work on the website, cd into the `website` folder and run
<code type="make/command" dir="website">make serve</code> to start a local
development server. The production site auto-updates on changes to the `main`
branch. The site hoster is [Netlify](https://www.netlify.com). Netlify
configuration is in [netlify.toml](netlify.toml).
