package godog

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	messages "github.com/cucumber/messages/go/v21"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/models"
	"github.com/cucumber/godog/internal/storage"
	"github.com/cucumber/godog/internal/utils"
)

var (
	errorInterface   = reflect.TypeOf((*error)(nil)).Elem()
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// more than one regex matched the step text
var ErrAmbiguous = fmt.Errorf("ambiguous step definition")

// ErrUndefined is returned in case if step definition was not found
var ErrUndefined = fmt.Errorf("step is undefined")

// ErrPending should be returned by step definition if
// step implementation is pending
var ErrPending = fmt.Errorf("step implementation is pending")

// ErrSkip should be returned by step definition or a hook if scenario and further steps are to be skipped.
var ErrSkip = fmt.Errorf("skipped")

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
	// StepAmbiguous indicates step text matches more than one step def
	StepAmbiguous = models.Ambiguous
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

type Attachment struct {
	Body      []byte
	FileName  string
	MediaType string
}

type attachmentKey struct{}

func Attach(ctx context.Context, attachments ...Attachment) context.Context {
	existing := Attachments(ctx)
	updated := append(existing, attachments...)
	return context.WithValue(ctx, attachmentKey{}, updated)
}

func Attachments(ctx context.Context) []Attachment {
	v := ctx.Value(attachmentKey{})

	if v == nil {
		return []Attachment{}
	}
	return v.([]Attachment)
}

func clearAttach(ctx context.Context) context.Context {
	return context.WithValue(ctx, attachmentKey{}, nil)
}

func pickleAttachments(ctx context.Context) []models.PickleAttachment {

	pickledAttachments := []models.PickleAttachment{}
	attachments := Attachments(ctx)

	for _, a := range attachments {
		pickledAttachments = append(pickledAttachments, models.PickleAttachment{
			Name:     a.FileName,
			Data:     a.Body,
			MimeType: a.MediaType,
		})
	}

	return pickledAttachments
}

func (s *suite) matchStep(step *messages.PickleStep) (*models.StepDefinition, error) {
	def, err := s.matchStepTextAndType(step.Text, step.Type)
	if err != nil {
		return nil, err
	}

	if def != nil && step.Argument != nil {
		def.Args = append(def.Args, step.Argument)
	}
	return def, nil
}

func (s *suite) runStep(ctx context.Context, pickle *Scenario, step *Step, scenarioErr error, isFirst, isLast bool) (rctx context.Context, err error) {
	var match *models.StepDefinition

	rctx = ctx

	// user multistep definitions may panic
	defer func() {
		if e := recover(); e != nil {
			pe, isErr := e.(error)
			switch {
			case isErr && errors.Is(pe, errStopNow):
				// FailNow or SkipNow called on dogTestingT, so clear the error to let the normal
				// below getTestingT(ctx).isFailed() call handle the reasons.
				err = nil
			case err != nil:
				err = &traceError{
					msg:   fmt.Sprintf("%s: %v", err.Error(), e),
					stack: callStack(),
				}
			default:
				err = &traceError{
					msg:   fmt.Sprintf("%v", e),
					stack: callStack(),
				}
			}
		}

		earlyReturn := scenarioErr != nil || errors.Is(err, ErrUndefined)

		// Check for any calls to Fail on dogT
		if err == nil {
			err = getTestingT(ctx).isFailed()
		}

		status := StepUndefined

		switch {
		case errors.Is(err, ErrAmbiguous):
			status = StepAmbiguous
		case errors.Is(err, ErrPending):
			status = StepPending
		case errors.Is(err, ErrSkip), err == nil && scenarioErr != nil:
			status = StepSkipped
		case errors.Is(err, ErrUndefined):
			status = StepUndefined
		case err != nil:
			status = StepFailed
		case err == nil && scenarioErr == nil:
			status = StepPassed
		}

		// Run after step handlers.
		rctx, err = s.runAfterStepHooks(ctx, step, status, err)

		// Trigger after scenario on failing or last step to attach possible hook error to step.
		if !s.shouldFail(scenarioErr) && (isLast || s.shouldFail(err)) {
			rctx, err = s.runAfterScenarioHooks(rctx, pickle, err)
		}

		// extract any accumulated attachments and clear them
		pickledAttachments := pickleAttachments(rctx)
		rctx = clearAttach(rctx)

		if earlyReturn {
			return
		}

		switch {
		case err == nil:
			sr := models.NewStepResult(models.Passed, pickle.Id, step.Id, match, pickledAttachments, nil)
			s.storage.MustInsertPickleStepResult(sr)
			s.fmt.Passed(pickle, step, match.GetInternalStepDefinition())
		case errors.Is(err, ErrPending):
			sr := models.NewStepResult(models.Pending, pickle.Id, step.Id, match, pickledAttachments, nil)
			s.storage.MustInsertPickleStepResult(sr)
			s.fmt.Pending(pickle, step, match.GetInternalStepDefinition())
		case errors.Is(err, ErrSkip):
			sr := models.NewStepResult(models.Skipped, pickle.Id, step.Id, match, pickledAttachments, nil)
			s.storage.MustInsertPickleStepResult(sr)
			s.fmt.Skipped(pickle, step, match.GetInternalStepDefinition())
		case errors.Is(err, ErrAmbiguous):
			sr := models.NewStepResult(models.Ambiguous, pickle.Id, step.Id, match, pickledAttachments, err)
			s.storage.MustInsertPickleStepResult(sr)
			s.fmt.Ambiguous(pickle, step, match.GetInternalStepDefinition(), err)
		default:
			sr := models.NewStepResult(models.Failed, pickle.Id, step.Id, match, pickledAttachments, err)
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

	var matchError error
	match, matchError = s.matchStep(step)

	s.storage.MustInsertStepDefintionMatch(step.AstNodeIds[0], match)
	s.fmt.Defined(pickle, step, match.GetInternalStepDefinition())

	if err != nil {
		pickledAttachments := pickleAttachments(ctx)
		ctx = clearAttach(ctx)

		sr := models.NewStepResult(models.Failed, pickle.Id, step.Id, match, pickledAttachments, nil)
		s.storage.MustInsertPickleStepResult(sr)
		return ctx, err
	}

	if matchError != nil {
		return ctx, matchError
	}

	if ctx, undef, err := s.maybeUndefined(ctx, step.Text, step.Argument, step.Type); err != nil {
		return ctx, err
	} else if len(undef) > 0 {
		if match != nil {
			match = &models.StepDefinition{
				StepDefinition: formatters.StepDefinition{
					Expr:    match.Expr,
					Handler: match.Handler,
					Keyword: match.Keyword,
				},
				Args:         match.Args,
				HandlerValue: match.HandlerValue,
				Nested:       match.Nested,
				Undefined:    undef,
			}
		}

		pickledAttachments := pickleAttachments(ctx)
		ctx = clearAttach(ctx)

		sr := models.NewStepResult(models.Undefined, pickle.Id, step.Id, match, pickledAttachments, nil)
		s.storage.MustInsertPickleStepResult(sr)

		s.fmt.Undefined(pickle, step, match.GetInternalStepDefinition())
		return ctx, ErrUndefined
	}

	if scenarioErr != nil {
		pickledAttachments := pickleAttachments(ctx)
		ctx = clearAttach(ctx)

		sr := models.NewStepResult(models.Skipped, pickle.Id, step.Id, match, pickledAttachments, nil)
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

func (s *suite) maybeUndefined(ctx context.Context, text string, arg interface{}, stepType messages.PickleStepType) (context.Context, []string, error) {
	var undefined []string
	step, err := s.matchStepTextAndType(text, stepType)
	if err != nil {
		return ctx, undefined, err
	}

	if nil == step {
		return ctx, []string{text}, nil
	}

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
		ctx, undef, err := s.maybeUndefined(ctx, next, nil, messages.PickleStepType_UNKNOWN)
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
		return ctx, fmt.Errorf("unexpected error, should have been godog.Steps: %T - %+v", result, result)
	}

	for _, text := range steps {
		def, err := s.matchStepTextAndType(text, messages.PickleStepType_UNKNOWN)
		if err != nil {
			return ctx, err
		}

		if def == nil {
			return ctx, ErrUndefined
		} else {
			ctx, err = s.runSubStep(ctx, text, def)
			if err != nil {
				return ctx, err
			}
		}
	}
	return ctx, nil
}

func (s *suite) runSubStep(ctx context.Context, text string, def *models.StepDefinition) (_ context.Context, err error) {
	st := &Step{}
	st.Text = text
	st.Type = messages.PickleStepType_ACTION

	defer func() {
		status := StepPassed

		switch {
		case errors.Is(err, ErrUndefined):
			status = StepUndefined
		case errors.Is(err, ErrPending):
			status = StepPending
		case err != nil:
			status = StepFailed
		}

		ctx, err = s.runAfterStepHooks(ctx, st, status, err)
	}()

	ctx, err = s.runBeforeStepHooks(ctx, st, nil)
	if err != nil {
		return ctx, fmt.Errorf("%s: %+v", text, err)
	}

	if ctx, err = s.maybeSubSteps(def.Run(ctx)); err != nil {
		return ctx, fmt.Errorf("%s: %+v", text, err)
	}

	return ctx, nil
}

func (s *suite) matchStepTextAndType(text string, stepType messages.PickleStepType) (*models.StepDefinition, error) {
	var first *models.StepDefinition
	matchingExpressions := make([]string, 0)

	for _, h := range s.steps {
		if m := h.Expr.FindStringSubmatch(text); len(m) > 0 {
			if !keywordMatches(h.Keyword, stepType) {
				continue
			}
			var args []interface{}
			for _, m := range m[1:] {
				args = append(args, m)
			}

			matchingExpressions = append(matchingExpressions, h.Expr.String())

			// since we need to assign arguments
			// better to copy the step definition
			match := &models.StepDefinition{
				StepDefinition: formatters.StepDefinition{
					Expr:    h.Expr,
					Handler: h.Handler,
					Keyword: h.Keyword,
				},
				Args:         args,
				HandlerValue: h.HandlerValue,
				Nested:       h.Nested,
			}

			if first == nil {
				first = match
			}
		}
	}

	if s.strict {
		if len(matchingExpressions) > 1 {
			errs := "\n        " + strings.Join(matchingExpressions, "\n        ")
			return nil, fmt.Errorf("%w, step text: %s\n    matches:%s", ErrAmbiguous, text, errs)
		}
	}

	return first, nil
}

func keywordMatches(k formatters.Keyword, stepType messages.PickleStepType) bool {
	if k == formatters.None {
		return true
	}
	switch stepType {
	case messages.PickleStepType_CONTEXT:
		return k == formatters.Given
	case messages.PickleStepType_ACTION:
		return k == formatters.When
	case messages.PickleStepType_OUTCOME:
		return k == formatters.Then
	default:
		return true
	}
}

func (s *suite) runSteps(ctx context.Context, pickle *Scenario, steps []*Step) (context.Context, error) {
	var (
		stepErr, scenarioErr error
	)

	for i, step := range steps {
		isLast := i == len(steps)-1
		isFirst := i == 0
		ctx, stepErr = s.runStep(ctx, pickle, step, scenarioErr, isFirst, isLast)
		if scenarioErr == nil || s.shouldFail(stepErr) {
			scenarioErr = stepErr
		}
	}

	return ctx, scenarioErr
}

func (s *suite) shouldFail(err error) bool {
	if err == nil || errors.Is(err, ErrSkip) {
		return false
	}

	if errors.Is(err, ErrUndefined) || errors.Is(err, ErrPending) {
		return s.strict
	}

	return true
}

func (s *suite) runPickle(pickle *messages.Pickle) (err error) {
	ctx := s.defaultContext
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

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

	dt := &testingT{
		name: pickle.Name,
	}
	ctx = setContextTestingT(ctx, dt)
	// scenario
	if s.testingT != nil {
		// Running scenario as a subtest.
		s.testingT.Run(pickle.Name, func(t *testing.T) {
			dt.t = t
			ctx, err = s.runSteps(ctx, pickle, pickle.Steps)
			if s.shouldFail(err) {
				t.Errorf("%+v", err)
			}
		})
	} else {
		ctx, err = s.runSteps(ctx, pickle, pickle.Steps)
	}

	// After scenario handlers are called in context of last evaluated step
	// so that error from handler can be added to step.

	return err
}
