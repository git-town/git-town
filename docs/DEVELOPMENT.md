# Developing the Git Town source code

This page provides guidance for contributing to the Git Town codebase. For a
comprehensive understanding of the architecture, refer to
[ARCHITECTURE.md](ARCHITECTURE.md).

## setup

1. install [Go](https://golang.org) version 1.21
2. install [Make](https://www.gnu.org/software/make)
   - Mac and Linux users have this out of the box
   - Windows users can install
     [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm) or
     run `choco install make` if [Chocolatey](https://chocolatey.org) is
     available.
3. run all CI tests locally: <code type="make/command" dir="..">make test</code>
4. faster smoke test during development: `make test-go`
5. install Git Town locally into `~/go/bin`:
   <code type="make/command" dir="..">make build</code>

## dependencies

Add an external Go dependency:

- run `go get [dependency]` inside the Git Town folder to register the
  dependency
- use the new dependency somewhere in the code
- run `go mod vendor` to vendor it
- run `go mod tidy` to clean up

Update an external Go module:

```
go get <path>
```

Update all external Go modules:

<a type="make/command" dir="..">

```
make update
```

</a>

## unit tests

Run unit tests for packages containing changes:

<a type="make/command" dir="..">

```
make unit
```

</a>

Run all unit tests no matter what has changed:

<a type="make/command" dir="..">

```
make unit-all
```

</a>

Run all unit tests with race detection:

<a type="make/command" dir="..">

```
make unit-race
```

</a>

Run an individual unit test:

```
go test src/cmd/root_test.go
go test src/cmd/root_test.go -v -run TestIsAcceptableGitVersion
```

## end-to-end tests

Run all end-to-end tests:

<a type="make/command" dir="..">

```
make cuke
```

</a>

Run all tests in the `features/append` folder or file:

```
go test -- features/append
```

To run individual Cucumber scenarios, add a `@this` flag to the scenario you
want to run. Example:

```cucumber
@this
Scenario: my awesome scenario
```

Then run only the scenarios that have a `@this` tag:

```
make cukethis
```

Certain tests require that the Git remote points to an actual GitHub, Gitea,
GitLab or Bitbucket address. This causes `git push` operations in this test to
also go to GitHub. To prevent this, set an environment variable
`GIT_TOWN_REMOTE` with the desired value of the `origin` remote, and Git Town
will use that value instead of what is configured in the repo.

If Cucumber tests produce garbled output on Windows, try running them inside Git
Bash. See [this issue](https://github.com/cucumber/godog/issues/129) for
details.

To pause an end-to-end test so that you have time to inspect the status of the
Git repository created by the test, add the step `And inspect the repo`. The
test runner will pause and print the path of the test workspace. You can `cd`
into that path in a separate terminal window and inspect the repos there. The
developer's repo is in the `repo` folder. The origin repo is in the `origin`
folder.

To see all commit SHAs of the repo, add the `And inspect the commits` step.

## inspecting variables

Inspect basic variables in a unit test:

```go
fmt.Printf("%#v\n", variable)
```

Inspect more complex variables:

```go
import "github.com/davecgh/go-spew/spew"

spew.Dump(variable)
```

- or -

```go
pretty.LDiff(t, var1, var2)
```

## debug end-to-end tests

To see the CLI output of the shell commands in a Cucumber test, as well as the
Git commands that the Git Town test suite runs under the hood, add a tag
`@debug` above the feature or scenario you want to debug:

```cucumber
@debug @this
Scenario: my awesome scenario
```

To see all Git commands that the test runner and the Git Town command run,
execute the Git Town command with the `--verbose` option. As an example, if the
step `When I run "git-town append new"` mysteriously fails, you could change it
to `When I run "git-town append new -v"`. Also add the tags `@debug @this` to
see the CLI output on the console.

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
doesn't enter. To find the hanging test you can do a binary search by executing
subsets of tests using `go test -- features/<path>` where path is either a
subfolder or file inside the "features" folder.

Alternatively, open `main_test.go`, change `Format` to `pretty` and
`Concurrency` to 1, and run the entire test suite. The detailed output will give
you hints at which test fails.

## run linters

Quick and efficient linter during development:

```
make lint
```

Run all linters reliably like on CI:

```
make lint-all
```

Auto-fix all fixable issues, including code formatting:

<a type="make/command" dir="..">

```
make fix
```

</a>

## debug the dialogs

Run `git town debug` to see the commands to manually test Git Town's dialogs.

## learn about the code and test architecture

See file [ARCHITECTURE.md](ARCHITECTURE.md).

## website

The source code for the [website](https://www.git-town.com) is in the
[website](../website) folder. This folder contains its own
[Makefile](../website/Makefile) for activities related to working on the
website.

To work on the website, cd into the `website` folder and run
<code type="make/command" dir="../website">make serve</code> to start a local
development server. The production site auto-updates on changes to the `main`
branch. The site hoster is [Netlify](https://www.netlify.com). Netlify
configuration is in [netlify.toml](../netlify.toml).
