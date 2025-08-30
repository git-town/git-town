package godog

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/builder"
	"github.com/cucumber/godog/internal/flags"
	"github.com/cucumber/godog/internal/models"
	messages "github.com/cucumber/messages/go/v21"
)

// GherkinDocument represents gherkin document.
type GherkinDocument = messages.GherkinDocument

// Scenario represents the executed scenario
type Scenario = messages.Pickle

// Step represents the executed step
type Step = messages.PickleStep

// Steps allows to nest steps
// instead of returning an error in step func
// it is possible to return combined steps:
//
//	func multistep(name string) godog.Steps {
//	  return godog.Steps{
//	    fmt.Sprintf(`an user named "%s"`, name),
//	    fmt.Sprintf(`user "%s" is authenticated`, name),
//	  }
//	}
//
// These steps will be matched and executed in
// sequential order. The first one which fails
// will result in main step failure.
type Steps []string

// StepDefinition is a registered step definition
// contains a StepHandler and regexp which
// is used to match a step. Args which
// were matched by last executed step
//
// This structure is passed to the formatter
// when step is matched and is either failed
// or successful
type StepDefinition = formatters.StepDefinition

// DocString represents the DocString argument made to a step definition
type DocString = messages.PickleDocString

// Table represents the Table argument made to a step definition
type Table = messages.PickleTable

// TestSuiteContext allows various contexts
// to register event handlers.
//
// When running a test suite, the instance of TestSuiteContext
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information
type TestSuiteContext struct {
	beforeSuiteHandlers []func()
	afterSuiteHandlers  []func()

	suite *suite
}

// BeforeSuite registers a function or method
// to be run once before suite runner.
//
// Use it to prepare the test suite for a spin.
// Connect and prepare database for instance...
func (ctx *TestSuiteContext) BeforeSuite(fn func()) {
	ctx.beforeSuiteHandlers = append(ctx.beforeSuiteHandlers, fn)
}

// AfterSuite registers a function or method
// to be run once after suite runner
func (ctx *TestSuiteContext) AfterSuite(fn func()) {
	ctx.afterSuiteHandlers = append(ctx.afterSuiteHandlers, fn)
}

// ScenarioContext allows registering scenario hooks.
func (ctx *TestSuiteContext) ScenarioContext() *ScenarioContext {
	return &ScenarioContext{
		suite: ctx.suite,
	}
}

// ScenarioContext allows various contexts
// to register steps and event handlers.
//
// When running a scenario, the instance of ScenarioContext
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information. Only step
// executions are catching panic error since it may
// be a context specific error.
type ScenarioContext struct {
	suite *suite
}

// StepContext allows registering step hooks.
type StepContext struct {
	suite *suite
}

// Before registers a function or method
// to be run before every scenario.
//
// It is a good practice to restore the default state
// before every scenario, so it would be isolated from
// any kind of state.
func (ctx ScenarioContext) Before(h BeforeScenarioHook) {
	ctx.suite.beforeScenarioHandlers = append(ctx.suite.beforeScenarioHandlers, h)
}

// BeforeScenarioHook defines a hook before scenario.
type BeforeScenarioHook func(ctx context.Context, sc *Scenario) (context.Context, error)

// After registers a function or method
// to be run after every scenario.
func (ctx ScenarioContext) After(h AfterScenarioHook) {
	ctx.suite.afterScenarioHandlers = append(ctx.suite.afterScenarioHandlers, h)
}

// AfterScenarioHook defines a hook after scenario.
type AfterScenarioHook func(ctx context.Context, sc *Scenario, err error) (context.Context, error)

// StepContext exposes StepContext of a scenario.
func (ctx ScenarioContext) StepContext() StepContext {
	return StepContext(ctx)
}

// Before registers a function or method
// to be run before every step.
func (ctx StepContext) Before(h BeforeStepHook) {
	ctx.suite.beforeStepHandlers = append(ctx.suite.beforeStepHandlers, h)
}

// BeforeStepHook defines a hook before step.
type BeforeStepHook func(ctx context.Context, st *Step) (context.Context, error)

// After registers a function or method
// to be run after every step.
//
// It may be convenient to return a different kind of error
// in order to print more state details which may help
// in case of step failure
//
// In some cases, for example when running a headless
// browser, to take a screenshot after failure.
func (ctx StepContext) After(h AfterStepHook) {
	ctx.suite.afterStepHandlers = append(ctx.suite.afterStepHandlers, h)
}

// AfterStepHook defines a hook after step.
type AfterStepHook func(ctx context.Context, st *Step, status StepResultStatus, err error) (context.Context, error)

// BeforeScenario registers a function or method
// to be run before every scenario.
//
// It is a good practice to restore the default state
// before every scenario, so it would be isolated from
// any kind of state.
//
// Deprecated: use Before.
func (ctx ScenarioContext) BeforeScenario(fn func(sc *Scenario)) {
	ctx.Before(func(ctx context.Context, sc *Scenario) (context.Context, error) {
		fn(sc)

		return ctx, nil
	})
}

// AfterScenario registers a function or method
// to be run after every scenario.
//
// Deprecated: use After.
func (ctx ScenarioContext) AfterScenario(fn func(sc *Scenario, err error)) {
	ctx.After(func(ctx context.Context, sc *Scenario, err error) (context.Context, error) {
		fn(sc, err)

		return ctx, nil
	})
}

// BeforeStep registers a function or method
// to be run before every step.
//
// Deprecated: use ScenarioContext.StepContext() and StepContext.Before.
func (ctx ScenarioContext) BeforeStep(fn func(st *Step)) {
	ctx.StepContext().Before(func(ctx context.Context, st *Step) (context.Context, error) {
		fn(st)

		return ctx, nil
	})
}

// AfterStep registers a function or method
// to be run after every step.
//
// It may be convenient to return a different kind of error
// in order to print more state details which may help
// in case of step failure
//
// In some cases, for example when running a headless
// browser, to take a screenshot after failure.
//
// Deprecated: use ScenarioContext.StepContext() and StepContext.After.
func (ctx ScenarioContext) AfterStep(fn func(st *Step, err error)) {
	ctx.StepContext().After(func(ctx context.Context, st *Step, status StepResultStatus, err error) (context.Context, error) {
		fn(st, err)

		return ctx, nil
	})
}

// Step allows to register a *StepDefinition in the
// Godog feature suite, the definition will be applied
// to all steps matching the given Regexp expr.
//
// It will panic if expr is not a valid regular
// expression or stepFunc is not a valid step
// handler.
//
// The expression can be of type: *regexp.Regexp, string or []byte
//
// The stepFunc may accept one or several arguments of type:
// - int, int8, int16, int32, int64
// - float32, float64
// - string
// - []byte
// - *godog.DocString
// - *godog.Table
//
// The stepFunc need to return either an error or []string for multistep
//
// Note that if there are two definitions which may match
// the same step, then only the first matched handler
// will be applied.
//
// If none of the *StepDefinition is matched, then
// ErrUndefined error will be returned when
// running steps.
func (ctx ScenarioContext) Step(expr, stepFunc interface{}) {
	ctx.stepWithKeyword(expr, stepFunc, formatters.None)
}

// Given functions identically to Step, but the *StepDefinition
// will only be matched if the step starts with "Given". "And"
// and "But" keywords copy the keyword of the last step for the
// purpose of matching.
func (ctx ScenarioContext) Given(expr, stepFunc interface{}) {
	ctx.stepWithKeyword(expr, stepFunc, formatters.Given)
}

// When functions identically to Step, but the *StepDefinition
// will only be matched if the step starts with "When". "And"
// and "But" keywords copy the keyword of the last step for the
// purpose of matching.
func (ctx ScenarioContext) When(expr, stepFunc interface{}) {
	ctx.stepWithKeyword(expr, stepFunc, formatters.When)
}

// Then functions identically to Step, but the *StepDefinition
// will only be matched if the step starts with "Then". "And"
// and "But" keywords copy the keyword of the last step for the
// purpose of matching.
func (ctx ScenarioContext) Then(expr, stepFunc interface{}) {
	ctx.stepWithKeyword(expr, stepFunc, formatters.Then)
}

func (ctx ScenarioContext) stepWithKeyword(expr interface{}, stepFunc interface{}, keyword formatters.Keyword) {
	var regex *regexp.Regexp

	// Validate the first input param is regex compatible
	switch t := expr.(type) {
	case *regexp.Regexp:
		regex = t
	case string:
		regex = regexp.MustCompile(t)
	case []byte:
		regex = regexp.MustCompile(string(t))
	default:
		panic(fmt.Sprintf("expecting expr to be a *regexp.Regexp or a string or []byte, got type: %T", expr))
	}

	// Validate that the handler is a function.
	handlerType := reflect.TypeOf(stepFunc)
	if handlerType.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected handler to be func, but got: %T", stepFunc))
	}

	// FIXME = Validate the handler function param types here so
	// that any errors are discovered early.
	// StepDefinition.Run defines the supported types but fails at run time not registration time

	// Validate the function's return types.
	helpPrefix := "expected handler to return one of error or context.Context or godog.Steps or (context.Context, error)"
	isNested := false

	numOut := handlerType.NumOut()
	switch numOut {
	case 0:
		// No return values.
	case 1:
		// One return value: should be error, Steps, or context.Context.
		outType := handlerType.Out(0)
		if outType == reflect.TypeOf(Steps{}) {
			isNested = true
		} else {
			if outType != errorInterface && outType != contextInterface {
				panic(fmt.Sprintf("%s, but got: %v", helpPrefix, outType))
			}
		}
	case 2:
		// Two return values: should be (context.Context, error).
		if handlerType.Out(0) != contextInterface || handlerType.Out(1) != errorInterface {
			panic(fmt.Sprintf("%s, but got: %v, %v", helpPrefix, handlerType.Out(0), handlerType.Out(1)))
		}
	default:
		// More than two return values.
		panic(fmt.Sprintf("expected handler to return either zero, one or two values, but it has: %d", numOut))
	}

	// Register the handler
	def := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: stepFunc,
			Expr:    regex,
			Keyword: keyword,
		},
		HandlerValue: reflect.ValueOf(stepFunc),
		Nested:       isNested,
	}

	// stash the step
	ctx.suite.steps = append(ctx.suite.steps, def)
}

// Build creates a test package like go test command at given target path.
// If there are no go files in tested directory, then
// it simply builds a godog executable to scan features.
//
// If there are go test files, it first builds a test
// package with standard go test command.
//
// Finally, it generates godog suite executable which
// registers exported godog contexts from the test files
// of tested package.
//
// Returns the path to generated executable
func Build(bin string) error {
	return builder.Build(bin)
}

type Feature = flags.Feature
