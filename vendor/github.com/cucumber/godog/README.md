[![#StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://vshymanskyy.github.io/StandWithUkraine)
[![Build Status](https://github.com/cucumber/godog/workflows/test/badge.svg)](https://github.com/cucumber/godog/actions?query=branch%main+workflow%3Atest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cucumber/godog)](https://pkg.go.dev/github.com/cucumber/godog)
[![codecov](https://codecov.io/gh/cucumber/godog/branch/master/graph/badge.svg)](https://codecov.io/gh/cucumber/godog)
[![pull requests](https://oselvar.com/api/badge?label=pull%20requests&csvUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fcucumber%2Foselvar-github-metrics%2Fmain%2Fdata%2Fcucumber%2Fgodog%2FpullRequests.csv)](https://oselvar.com/github/cucumber/oselvar-github-metrics/main/cucumber/godog)
[![issues](https://oselvar.com/api/badge?label=issues&csvUrl=https%3A%2F%2Fraw.githubusercontent.com%2Fcucumber%2Foselvar-github-metrics%2Fmain%2Fdata%2Fcucumber%2Fgodog%2Fissues.csv)](https://oselvar.com/github/cucumber/oselvar-github-metrics/main/cucumber/godog)

# Godog

<p align="center"><img src="logo.png" alt="Godog logo" style="width:250px;" /></p>

**The API is likely to change a few times before we reach 1.0.0**

Please read the full README, you may find it very useful. And do not forget to peek into the [Release Notes](https://github.com/cucumber/godog/blob/master/release-notes) and the [CHANGELOG](https://github.com/cucumber/godog/blob/master/CHANGELOG.md) from time to time.

Package godog is the official Cucumber BDD framework for Golang, it merges specification and test documentation into one cohesive whole, using Gherkin formatted scenarios in the format of Given, When, Then.

The project was inspired by [behat][behat] and [cucumber][cucumber].

## Why Godog/Cucumber

### A single source of truth

Godog merges specification and test documentation into one cohesive whole.

### Living documentation

Because they're automatically tested by Godog, your specifications are
always bang up-to-date.

### Focus on the customer

Business and IT don't always understand each other. Godog's executable specifications encourage closer collaboration, helping teams keep the business goal in mind at all times.

### Less rework

When automated testing is this much fun, teams can easily protect themselves from costly regressions.

### Read more
- [Behaviour-Driven Development](https://cucumber.io/docs/bdd/)
- [Gherkin Reference](https://cucumber.io/docs/gherkin/reference/)

## Contributions

Godog is a community driven Open Source Project within the Cucumber organization. We [welcome contributions from everyone](https://cucumber.io/blog/open-source/tackling-structural-racism-(and-sexism)-in-open-so/), and we're ready to support you if you have the enthusiasm to contribute.

See the [contributing guide] for more detail on how to get started.

## Getting help

We have a [community Slack] where you can chat with other users, developers, and BDD practitioners.

Here are some useful channels to try:

- [#help-godog](https://cucumberbdd.slack.com/archives/CTNL1JCVA) - General Godog Adoption Help
- [#committers-go](https://cucumberbdd.slack.com/archives/CA5NJPDJ4) - Golang focused Cucumber Contributors
- [#committers](https://cucumberbdd.slack.com/archives/C62D0FK0E) - General Cucumber Contributors

## Examples

You can find a few examples [here](/_examples).

**Note** that if you want to execute any of the examples and have the Git repository checked out in the `$GOPATH`, you need to use: `GO111MODULE=off`. [Issue](https://github.com/cucumber/godog/issues/344) for reference.

### Godogs

The following example can be [found here](/_examples/godogs).

#### Step 1 - Setup a go module

Given we create a new go module **godogs** in your normal go workspace. - `mkdir godogs`

From now on, this is our work directory - `cd godogs`

Initiate the go module - `go mod init godogs`

#### Step 2 - Install godog

Install the `godog` binary - `go install github.com/cucumber/godog/cmd/godog@latest`.

#### Step 3 - Create gherkin feature

Imagine we have a **godog cart** to serve godogs for lunch.

First of all, we describe our feature in plain text - `vim features/godogs.feature`.

``` gherkin
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
```

#### Step 4 - Create godog step definitions

**NOTE:** same as **go test** godog respects package level isolation. All your step definitions should be in your tested package root directory. In this case: **godogs**.

If we run godog inside the module: - `godog run`

You should see that the steps are undefined:
```
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12          # features/godogs.feature:6
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining

1 scenarios (1 undefined)
3 steps (3 undefined)
220.129µs

You can implement step definitions for undefined steps with these snippets:

func iEat(arg1 int) error {
        return godog.ErrPending
}

func thereAreGodogs(arg1 int) error {
        return godog.ErrPending
}

func thereShouldBeRemaining(arg1 int) error {
        return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
        ctx.Step(`^I eat (\d+)$`, iEat)
        ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
        ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
```

Create and copy the step definitions into a new file - `vim godogs_test.go`
``` go
package main

import "github.com/cucumber/godog"

func iEat(arg1 int) error {
        return godog.ErrPending
}

func thereAreGodogs(arg1 int) error {
        return godog.ErrPending
}

func thereShouldBeRemaining(arg1 int) error {
        return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
        ctx.Step(`^I eat (\d+)$`, iEat)
        ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
        ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
```

Our module should now look like this:
```
godogs
- features
  - godogs.feature
- go.mod
- go.sum
- godogs_test.go
```

Run godog again - `godog run`

You should now see that the scenario is pending with one step pending and two steps skipped:
```
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12          # features/godogs.feature:6
    Given there are 12 godogs        # godogs_test.go:10 -> thereAreGodogs
      TODO: write pending definition
    When I eat 5                     # godogs_test.go:6 -> iEat
    Then there should be 7 remaining # godogs_test.go:14 -> thereShouldBeRemaining

1 scenarios (1 pending)
3 steps (1 pending, 2 skipped)
282.123µs
```

You may change **return godog.ErrPending** to **return nil** in the three step definitions and the scenario will pass successfully.

Also, you may omit error return if your step does not fail.

```go
func iEat(arg1 int) {
	// Eat arg1.
}
```

#### Step 5 - Create the main program to test

We only need a number of **godogs** for now. Lets keep it simple.

Create and copy the code into a new file - `vim godogs.go`
```go
package main

// Godogs available to eat
var Godogs int

func main() { /* usual main func */ }
```

Our module should now look like this:
```
godogs
- features
  - godogs.feature
- go.mod
- go.sum
- godogs.go
- godogs_test.go
```

#### Step 6 - Add some logic to the step definitions

Now lets implement our step definitions to test our feature requirements:

Replace the contents of `godogs_test.go` with the code below - `vim godogs_test.go`

```go
package main

import (
  "context"
  "errors"
  "fmt"
  "testing"

  "github.com/cucumber/godog"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

func thereAreGodogs(ctx context.Context, available int) (context.Context, error) {
  return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func iEat(ctx context.Context, num int) (context.Context, error) {
  available, ok := ctx.Value(godogsCtxKey{}).(int)
  if !ok {
    return ctx, errors.New("there are no godogs available")
  }

  if available < num {
    return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
  }

  available -= num

  return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
  available, ok := ctx.Value(godogsCtxKey{}).(int)
  if !ok {
    return errors.New("there are no godogs available")
  }

  if available != remaining {
    return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, available)
  }

  return nil
}

func TestFeatures(t *testing.T) {
  suite := godog.TestSuite{
    ScenarioInitializer: InitializeScenario,
    Options: &godog.Options{
      Format:   "pretty",
      Paths:    []string{"features"},
      TestingT: t, // Testing instance that will run subtests.
    },
  }

  if suite.Run() != 0 {
    t.Fatal("non-zero status returned, failed to run feature tests")
  }
}

func InitializeScenario(sc *godog.ScenarioContext) {
  sc.Step(`^there are (\d+) godogs$`, thereAreGodogs)
  sc.Step(`^I eat (\d+)$`, iEat)
  sc.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}

```

In this example we are using `context.Context` to pass the state between the steps. 
Every scenario starts with an empty context and then steps and hooks can add relevant information to it.
Instrumented context is chained through the steps and hooks and is safe to use when multiple scenarios are running concurrently.

When you run godog again with `go test -v godogs_test.go` or with a CLI `godog run`.

You should see a passing run:
```
=== RUN   TestFeatures
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs
=== RUN   TestFeatures/Eat_5_out_of_12

  Scenario: Eat 5 out of 12          # features/godogs.feature:6
    Given there are 12 godogs        # godogs_test.go:14 -> command-line-arguments.thereAreGodogs
    When I eat 5                     # godogs_test.go:18 -> command-line-arguments.iEat
    Then there should be 7 remaining # godogs_test.go:33 -> command-line-arguments.thereShouldBeRemaining

1 scenarios (1 passed)
3 steps (3 passed)
275.333µs
--- PASS: TestFeatures (0.00s)
    --- PASS: TestFeatures/Eat_5_out_of_12 (0.00s)
PASS
ok      command-line-arguments  0.130s
```

You may hook to `ScenarioContext` **Before** event in order to reset or pre-seed the application state before each scenario. 
You may hook into more events, like `sc.StepContext()` **After** to print all state in case of an error. 
Or **BeforeSuite** to prepare a database.

By now, you should have figured out, how to use **godog**. Another advice is to make steps orthogonal, small and simple to read for a user. Whether the user is a dumb website user or an API developer, who may understand a little more technical context - it should target that user.

When steps are orthogonal and small, you can combine them just like you do with Unix tools. Look how to simplify or remove ones, which can be composed.

`TestFeatures` acts as a regular Go test, so you can leverage your IDE facilities to run and debug it.

## Code of Conduct

Everyone interacting in this codebase and issue tracker is expected to follow the Cucumber [code of conduct](https://github.com/cucumber/cucumber/blob/master/CODE_OF_CONDUCT.md).

## References and Tutorials

- [cucumber-html-reporter](https://github.com/gkushang/cucumber-html-reporter),
  may be used in order to generate **html** reports together with **cucumber** output formatter. See the [following docker image](https://github.com/myie/cucumber-html-reporter) for usage details.
- [how to use godog by semaphoreci](https://semaphoreci.com/community/tutorials/how-to-use-godog-for-behavior-driven-development-in-go)
- see [examples](https://github.com/cucumber/godog/tree/master/_examples)
- see extension [AssistDog](https://github.com/hellomd/assistdog),
  which may have useful **gherkin.DataTable** transformations or comparison methods for assertions.

## Documentation

See [pkg documentation][godoc] for general API details.
See **[Circle Config](/.circleci/config.yml)** for supported **go** versions.
See `godog -h` for general command options.

See implementation examples:

- [rest API server](/_examples/api)
- [rest API with Database](/_examples/db)
- [godogs](/_examples/godogs)

## FAQ

### Running Godog with go test

You may integrate running **godog** in your **go test** command. 

#### Subtests of *testing.T

You can run test suite using go [Subtests](https://pkg.go.dev/testing#hdr-Subtests_and_Sub_benchmarks).
In this case it is not necessary to have **godog** command installed. See the following example.

```go
package main_test

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
  suite := godog.TestSuite{
    ScenarioInitializer: func(s *godog.ScenarioContext) {
      // Add step definitions here.
    },
    Options: &godog.Options{
      Format:   "pretty",
      Paths:    []string{"features"},
      TestingT: t, // Testing instance that will run subtests.
    },
  }

  if suite.Run() != 0 {
    t.Fatal("non-zero status returned, failed to run feature tests")
  }
}
```

Then you can run suite.
```
go test -test.v -test.run ^TestFeatures$
```

Or a particular scenario.
```
go test -test.v -test.run ^TestFeatures$/^my_scenario$
```

#### TestMain

You can run test suite using go [TestMain](https://golang.org/pkg/testing/#hdr-Main) func available since **go 1.4**. 
In this case it is not necessary to have **godog** command installed. See the following examples.

The following example binds **godog** flags with specified prefix `godog` in order to prevent flag collisions.

```go
package main

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag" // godog v0.11.0 and later
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", pflag.CommandLine, &opts) // godog v0.10.0 and earlier
	godog.BindCommandLineFlags("godog.", &opts)        // godog v0.11.0 and later
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
```

Then you may run tests with by specifying flags in order to filter features.

```
go test -v --godog.random --godog.tags=wip
go test -v --godog.format=pretty --godog.random -race -coverprofile=coverage.txt -covermode=atomic
```

The following example does not bind godog flags, instead manually configuring needed options.

```go
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
```

You can even go one step further and reuse **go test** flags, like **verbose** mode in order to switch godog **format**. See the following example:

```go
func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format: format,
		Paths:     []string{"features"},
	}

	status := godog.TestSuite{
		Name: "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
```

Now when running `go test -v` it will use **pretty** format.

### Tags

If you want to filter scenarios by tags, you can use the `-t=<expression>` or `--tags=<expression>` where `<expression>` is one of the following:

- `@wip` - run all scenarios with wip tag
- `~@wip` - exclude all scenarios with wip tag
- `@wip && ~@new` - run wip scenarios, but exclude new
- `@wip,@undone` - run wip or undone scenarios

### Using assertion packages like testify with Godog
A more extensive example can be [found here](/_examples/assert-godogs).

```go
func thereShouldBeRemaining(remaining int) error {
	return assertExpectedAndActual(
		assert.Equal, Godogs, remaining,
		"Expected %d godogs to be remaining, but there is %d", remaining, Godogs,
	)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}
```

## CLI Mode

Another way to use `godog` is to run it in CLI mode.

In this mode `godog` CLI will use `go` under the hood to compile and run your test suite.

**Godog** does not intervene with the standard **go test** command behavior. You can leverage both frameworks to functionally test your application while maintaining all test related source code in **_test.go** files.

**Godog** acts similar compared to **go test** command, by using go compiler and linker tool in order to produce test executable. Godog contexts need to be exported the same way as **Test** functions for go tests. Note, that if you use **godog** command tool, it will use `go` executable to determine compiler and linker.

### Install
```
go install github.com/cucumber/godog/cmd/godog@latest
```
Adding `@v0.12.0` will install v0.12.0 specifically instead of master.

With `go` version prior to 1.17, use `go get github.com/cucumber/godog/cmd/godog@v0.12.0`.
Running `within the $GOPATH`, you would also need to set `GO111MODULE=on`, like this:
```
GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.12.0
```

### Configure common options for godog CLI

There are no global options or configuration files. Alias your common or project based commands: `alias godog-wip="godog --format=progress --tags=@wip"`

## Concurrency

When concurrency is configured in options, godog will execute the scenarios concurrently, which is supported by all supplied formatters.

In order to support concurrency well, you should reset the state and isolate each scenario. They should not share any state. It is suggested to run the suite concurrently in order to make sure there is no state corruption or race conditions in the application.

It is also useful to randomize the order of scenario execution, which you can now do with `--random` command option or `godog.Options.Randomize` setting.

### Building your own custom formatter
A simple example can be [found here](/_examples/custom-formatter).

## License
**Godog** and **Gherkin** are licensed under the [MIT][license] and developed as a part of the [cucumber project][cucumber]

[godoc]: https://pkg.go.dev/github.com/cucumber/godog "Documentation on godog"
[golang]: https://golang.org/  "GO programming language"
[behat]: http://docs.behat.org/ "Behavior driven development framework for PHP"
[cucumber]: https://cucumber.io/ "Behavior driven development framework"
[license]: https://en.wikipedia.org/wiki/MIT_License "The MIT license"
[contributing guide]: https://github.com/cucumber/godog/blob/main/CONTRIBUTING.md
[community Slack]: https://cucumber.io/community#slack
