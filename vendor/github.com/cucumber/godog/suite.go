package godog

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/models"
	"github.com/cucumber/godog/internal/storage"
	"github.com/cucumber/godog/internal/utils"
)

var (
	errorInterface   = reflect.TypeOf((*error)(nil)).Elem()
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// ErrUndefined is returned in case if step definition was not found
var ErrUndefined = fmt.Errorf("step is undefined")

// ErrPending should be returned by step definition if
// step implementation is pending
var ErrPending = fmt.Errorf("step implementation is pending")

// StepResultStatus describes step result.
type StepResultStatus = models.StepResultStatus

const (
	// StepPassed indicates step that passed.
	StepPassed StepResultStatus = models.Passed
	// StepFailed indicates step that failed.
	StepFailed = models.Failed
	// StepSkipped indicates step that was skipped.
	StepSkipped = models.Skipped
	// StepUndefined indicates undefined step.
	StepUndefined = models.Undefined
	// StepPending indicates step with pending implementation.
	StepPending = models.Pending
)

type suite struct {
	steps []*models.StepDefinition

	fmt     Formatter
	storage *storage.Storage

	failed        bool
	randomSeed    int64
	stopOnFailure bool
	strict        bool

	defaultContext context.Context
	testingT       *testing.T

	// suite event handlers
	beforeScenarioHandlers []BeforeScenarioHook
	beforeStepHandlers     []BeforeStepHook
	afterStepHandlers      []AfterStepHook
	afterScenarioHandlers  []AfterScenarioHook
}

func (s *suite) matchStep(step *messages.PickleStep) *models.StepDefinition {
	def := s.matchStepText(step.Text)
	if def != nil && step.Argument != nil {
		def.Args = append(def.Args, step.Argument)
	}
	return def
}

func (s *suite) runStep(ctx context.Context, pickle *Scenario, step *Step, prevStepErr error, isFirst, isLast bool) (rctx context.Context, err error) {
	var (
		match *models.StepDefinition
		sr    = models.PickleStepResult{Status: models.Undefined}
	)

	rctx = ctx

	// user multistep definitions may panic
	defer func() {
		if e := recover(); e != nil {
			err = &traceError{
				msg:   fmt.Sprintf("%v", e),
				stack: callStack(),
			}
		}

		earlyReturn := prevStepErr != nil || err == ErrUndefined

		if !earlyReturn {
			sr = models.NewStepResult(pickle.Id, step.Id, match)
		}

		// Run after step handlers.
		rctx, err = s.runAfterStepHooks(ctx, step, sr.Status, err)

		// Trigger after scenario on failing or last step to attach possible hook error to step.
		if isLast || (sr.Status != StepSkipped && sr.Status != StepUndefined && err != nil) {
			rctx, err = s.runAfterScenarioHooks(rctx, pickle, err)
		}

		if earlyReturn {
			return
		}

		switch err {
		case nil:
			sr.Status = models.Passed
			s.storage.MustInsertPickleStepResult(sr)

			s.fmt.Passed(pickle, step, match.GetInternalStepDefinition())
		case ErrPending:
			sr.Status = models.Pending
			s.storage.MustInsertPickleStepResult(sr)

			s.fmt.Pending(pickle, step, match.GetInternalStepDefinition())
		default:
			sr.Status = models.Failed
			sr.Err = err
			s.storage.MustInsertPickleStepResult(sr)

			s.fmt.Failed(pickle, step, match.GetInternalStepDefinition(), err)
		}
	}()

	// run before scenario handlers
	if isFirst {
		ctx, err = s.runBeforeScenarioHooks(ctx, pickle)
	}

	// run before step handlers
	ctx, err = s.runBeforeStepHooks(ctx, step, err)

	match = s.matchStep(step)
	s.storage.MustInsertStepDefintionMatch(step.AstNodeIds[0], match)
	s.fmt.Defined(pickle, step, match.GetInternalStepDefinition())

	if err != nil {
		sr = models.NewStepResult(pickle.Id, step.Id, match)
		sr.Status = models.Failed
		s.storage.MustInsertPickleStepResult(sr)

		return ctx, err
	}

	if ctx, undef, err := s.maybeUndefined(ctx, step.Text, step.Argument); err != nil {
		return ctx, err
	} else if len(undef) > 0 {
		if match != nil {
			match = &models.StepDefinition{
				StepDefinition: formatters.StepDefinition{
					Expr:    match.Expr,
					Handler: match.Handler,
				},
				Args:         match.Args,
				HandlerValue: match.HandlerValue,
				Nested:       match.Nested,
				Undefined:    undef,
			}
		}

		sr = models.NewStepResult(pickle.Id, step.Id, match)
		sr.Status = models.Undefined
		s.storage.MustInsertPickleStepResult(sr)

		s.fmt.Undefined(pickle, step, match.GetInternalStepDefinition())
		return ctx, ErrUndefined
	}

	if prevStepErr != nil {
		sr = models.NewStepResult(pickle.Id, step.Id, match)
		sr.Status = models.Skipped
		s.storage.MustInsertPickleStepResult(sr)

		s.fmt.Skipped(pickle, step, match.GetInternalStepDefinition())
		return ctx, nil
	}

	ctx, err = s.maybeSubSteps(match.Run(ctx))

	return ctx, err
}

func (s *suite) runBeforeStepHooks(ctx context.Context, step *Step, err error) (context.Context, error) {
	hooksFailed := false

	for _, f := range s.beforeStepHandlers {
		hctx, herr := f(ctx, step)
		if herr != nil {
			hooksFailed = true

			if err == nil {
				err = herr
			} else {
				err = fmt.Errorf("%v, %w", herr, err)
			}
		}

		if hctx != nil {
			ctx = hctx
		}
	}

	if hooksFailed {
		err = fmt.Errorf("before step hook failed: %w", err)
	}

	return ctx, err
}

func (s *suite) runAfterStepHooks(ctx context.Context, step *Step, status StepResultStatus, err error) (context.Context, error) {
	for _, f := range s.afterStepHandlers {
		hctx, herr := f(ctx, step, status, err)

		// Adding hook error to resulting error without breaking hooks loop.
		if herr != nil {
			if err == nil {
				err = herr
			} else {
				err = fmt.Errorf("%v, %w", herr, err)
			}
		}

		if hctx != nil {
			ctx = hctx
		}
	}

	return ctx, err
}

func (s *suite) runBeforeScenarioHooks(ctx context.Context, pickle *messages.Pickle) (context.Context, error) {
	var err error

	// run before scenario handlers
	for _, f := range s.beforeScenarioHandlers {
		hctx, herr := f(ctx, pickle)
		if herr != nil {
			if err == nil {
				err = herr
			} else {
				err = fmt.Errorf("%v, %w", herr, err)
			}
		}

		if hctx != nil {
			ctx = hctx
		}
	}

	if err != nil {
		err = fmt.Errorf("before scenario hook failed: %w", err)
	}

	return ctx, err
}

func (s *suite) runAfterScenarioHooks(ctx context.Context, pickle *messages.Pickle, lastStepErr error) (context.Context, error) {
	err := lastStepErr

	hooksFailed := false
	isStepErr := true

	// run after scenario handlers
	for _, f := range s.afterScenarioHandlers {
		hctx, herr := f(ctx, pickle, err)

		// Adding hook error to resulting error without breaking hooks loop.
		if herr != nil {
			hooksFailed = true

			if err == nil {
				isStepErr = false
				err = herr
			} else {
				if isStepErr {
					err = fmt.Errorf("step error: %w", err)
					isStepErr = false
				}
				err = fmt.Errorf("%v, %w", herr, err)
			}
		}

		if hctx != nil {
			ctx = hctx
		}
	}

	if hooksFailed {
		err = fmt.Errorf("after scenario hook failed: %w", err)
	}

	return ctx, err
}

func (s *suite) maybeUndefined(ctx context.Context, text string, arg interface{}) (context.Context, []string, error) {
	step := s.matchStepText(text)
	if nil == step {
		return ctx, []string{text}, nil
	}

	var undefined []string
	if !step.Nested {
		return ctx, undefined, nil
	}

	if arg != nil {
		step.Args = append(step.Args, arg)
	}

	ctx, steps := step.Run(ctx)

	for _, next := range steps.(Steps) {
		lines := strings.Split(next, "\n")
		// @TODO: we cannot currently parse table or content body from nested steps
		if len(lines) > 1 {
			return ctx, undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		if len(lines[0]) > 0 && lines[0][len(lines[0])-1] == ':' {
			return ctx, undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		ctx, undef, err := s.maybeUndefined(ctx, next, nil)
		if err != nil {
			return ctx, undefined, err
		}
		undefined = append(undefined, undef...)
	}
	return ctx, undefined, nil
}

func (s *suite) maybeSubSteps(ctx context.Context, result interface{}) (context.Context, error) {
	if nil == result {
		return ctx, nil
	}

	if err, ok := result.(error); ok {
		return ctx, err
	}

	steps, ok := result.(Steps)
	if !ok {
		return ctx, fmt.Errorf("unexpected error, should have been []string: %T - %+v", result, result)
	}

	var err error

	for _, text := range steps {
		if def := s.matchStepText(text); def == nil {
			return ctx, ErrUndefined
		} else if ctx, err = s.maybeSubSteps(def.Run(ctx)); err != nil {
			return ctx, fmt.Errorf("%s: %+v", text, err)
		}
	}
	return ctx, nil
}

func (s *suite) matchStepText(text string) *models.StepDefinition {
	for _, h := range s.steps {
		if m := h.Expr.FindStringSubmatch(text); len(m) > 0 {
			var args []interface{}
			for _, m := range m[1:] {
				args = append(args, m)
			}

			// since we need to assign arguments
			// better to copy the step definition
			return &models.StepDefinition{
				StepDefinition: formatters.StepDefinition{
					Expr:    h.Expr,
					Handler: h.Handler,
				},
				Args:         args,
				HandlerValue: h.HandlerValue,
				Nested:       h.Nested,
			}
		}
	}
	return nil
}

func (s *suite) runSteps(ctx context.Context, pickle *Scenario, steps []*Step) (context.Context, error) {
	var (
		stepErr, err error
	)

	for i, step := range steps {
		isLast := i == len(steps)-1
		isFirst := i == 0
		ctx, stepErr = s.runStep(ctx, pickle, step, err, isFirst, isLast)
		switch stepErr {
		case ErrUndefined:
			// do not overwrite failed error
			if err == ErrUndefined || err == nil {
				err = stepErr
			}
		case ErrPending:
			err = stepErr
		case nil:
		default:
			err = stepErr
		}
	}

	return ctx, err
}

func (s *suite) shouldFail(err error) bool {
	if err == nil {
		return false
	}

	if err == ErrUndefined || err == ErrPending {
		return s.strict
	}

	return true
}

func (s *suite) runPickle(pickle *messages.Pickle) (err error) {
	ctx := s.defaultContext
	if ctx == nil {
		ctx = context.Background()
	}

	if len(pickle.Steps) == 0 {
		pr := models.PickleResult{PickleID: pickle.Id, StartedAt: utils.TimeNowFunc()}
		s.storage.MustInsertPickleResult(pr)

		s.fmt.Pickle(pickle)
		return ErrUndefined
	}

	// Before scenario hooks are called in context of first evaluated step
	// so that error from handler can be added to step.

	pr := models.PickleResult{PickleID: pickle.Id, StartedAt: utils.TimeNowFunc()}
	s.storage.MustInsertPickleResult(pr)

	s.fmt.Pickle(pickle)

	// scenario
	if s.testingT != nil {
		// Running scenario as a subtest.
		s.testingT.Run(pickle.Name, func(t *testing.T) {
			ctx, err = s.runSteps(ctx, pickle, pickle.Steps)
			if err != nil {
				t.Error(err)
			}
		})
	} else {
		ctx, err = s.runSteps(ctx, pickle, pickle.Steps)
	}

	// After scenario handlers are called in context of last evaluated step
	// so that error from handler can be added to step.

	return err
}
