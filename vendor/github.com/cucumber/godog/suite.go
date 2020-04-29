package godog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cucumber/gherkin-go/v11"
	"github.com/cucumber/messages-go/v10"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var typeOfBytes = reflect.TypeOf([]byte(nil))

type feature struct {
	*messages.GherkinDocument
	pickles       []*messages.Pickle
	pickleResults []*pickleResult

	time    time.Time
	Content []byte `json:"-"`
	Path    string `json:"path"`
	order   int
}

func (f feature) findScenario(astScenarioID string) *messages.GherkinDocument_Feature_Scenario {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil && sc.Id == astScenarioID {
			return sc
		}
	}

	return nil
}

func (f feature) findBackground(astScenarioID string) *messages.GherkinDocument_Feature_Background {
	var bg *messages.GherkinDocument_Feature_Background

	for _, child := range f.GherkinDocument.Feature.Children {
		if tmp := child.GetBackground(); tmp != nil {
			bg = tmp
		}

		if sc := child.GetScenario(); sc != nil && sc.Id == astScenarioID {
			return bg
		}
	}

	return nil
}

func (f feature) findExample(exampleAstID string) (*messages.GherkinDocument_Feature_Scenario_Examples, *messages.GherkinDocument_Feature_TableRow) {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil {
			for _, example := range sc.Examples {
				for _, row := range example.TableBody {
					if row.Id == exampleAstID {
						return example, row
					}
				}
			}
		}
	}

	return nil, nil
}

func (f feature) findStep(astStepID string) *messages.GherkinDocument_Feature_Step {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.GetScenario(); sc != nil {
			for _, step := range sc.GetSteps() {
				if step.Id == astStepID {
					return step
				}
			}
		}

		if bg := child.GetBackground(); bg != nil {
			for _, step := range bg.GetSteps() {
				if step.Id == astStepID {
					return step
				}
			}
		}
	}

	return nil
}

func (f feature) startedAt() time.Time {
	return f.time
}

func (f feature) finishedAt() time.Time {
	if len(f.pickleResults) == 0 {
		return f.startedAt()
	}

	return f.pickleResults[len(f.pickleResults)-1].finishedAt()
}

func (f feature) appendStepResult(s *stepResult) {
	pickles := f.pickleResults[len(f.pickleResults)-1]
	pickles.stepResults = append(pickles.stepResults, s)
}

func (f feature) lastPickleResult() *pickleResult {
	return f.pickleResults[len(f.pickleResults)-1]
}

func (f feature) lastStepResult() *stepResult {
	last := f.lastPickleResult()
	return last.stepResults[len(last.stepResults)-1]
}

type sortByName []*feature

func (s sortByName) Len() int           { return len(s) }
func (s sortByName) Less(i, j int) bool { return s[i].Feature.Name < s[j].Feature.Name }
func (s sortByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type pickleResult struct {
	Name        string
	time        time.Time
	stepResults []*stepResult
}

func (s pickleResult) startedAt() time.Time {
	return s.time
}

func (s pickleResult) finishedAt() time.Time {
	if len(s.stepResults) == 0 {
		return s.startedAt()
	}

	return s.stepResults[len(s.stepResults)-1].time
}

// ErrUndefined is returned in case if step definition was not found
var ErrUndefined = fmt.Errorf("step is undefined")

// ErrPending should be returned by step definition if
// step implementation is pending
var ErrPending = fmt.Errorf("step implementation is pending")

// Suite allows various contexts
// to register steps and event handlers.
//
// When running a test suite, the instance of Suite
// is passed to all functions (contexts), which
// have it as a first and only argument.
//
// Note that all event hooks does not catch panic errors
// in order to have a trace information. Only step
// executions are catching panic error since it may
// be a context specific error.
type Suite struct {
	steps    []*StepDefinition
	features []*feature
	fmt      Formatter

	failed        bool
	randomSeed    int64
	stopOnFailure bool
	strict        bool

	// suite event handlers
	beforeSuiteHandlers    []func()
	beforeFeatureHandlers  []func(*messages.GherkinDocument)
	beforeScenarioHandlers []func(*messages.Pickle)
	beforeStepHandlers     []func(*messages.Pickle_PickleStep)
	afterStepHandlers      []func(*messages.Pickle_PickleStep, error)
	afterScenarioHandlers  []func(*messages.Pickle, error)
	afterFeatureHandlers   []func(*messages.GherkinDocument)
	afterSuiteHandlers     []func()
}

// Step allows to register a *StepDefinition in Godog
// feature suite, the definition will be applied
// to all steps matching the given Regexp expr.
//
// It will panic if expr is not a valid regular
// expression or stepFunc is not a valid step
// handler.
//
// Note that if there are two definitions which may match
// the same step, then only the first matched handler
// will be applied.
//
// If none of the *StepDefinition is matched, then
// ErrUndefined error will be returned when
// running steps.
func (s *Suite) Step(expr interface{}, stepFunc interface{}) {
	var regex *regexp.Regexp

	switch t := expr.(type) {
	case *regexp.Regexp:
		regex = t
	case string:
		regex = regexp.MustCompile(t)
	case []byte:
		regex = regexp.MustCompile(string(t))
	default:
		panic(fmt.Sprintf("expecting expr to be a *regexp.Regexp or a string, got type: %T", expr))
	}

	v := reflect.ValueOf(stepFunc)
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected handler to be func, but got: %T", stepFunc))
	}

	if typ.NumOut() != 1 {
		panic(fmt.Sprintf("expected handler to return only one value, but it has: %d", typ.NumOut()))
	}

	def := &StepDefinition{
		Handler: stepFunc,
		Expr:    regex,
		hv:      v,
	}

	typ = typ.Out(0)
	switch typ.Kind() {
	case reflect.Interface:
		if !typ.Implements(errorInterface) {
			panic(fmt.Sprintf("expected handler to return an error, but got: %s", typ.Kind()))
		}
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.String {
			panic(fmt.Sprintf("expected handler to return []string for multistep, but got: []%s", typ.Kind()))
		}
		def.nested = true
	default:
		panic(fmt.Sprintf("expected handler to return an error or []string, but got: %s", typ.Kind()))
	}

	s.steps = append(s.steps, def)
}

// BeforeSuite registers a function or method
// to be run once before suite runner.
//
// Use it to prepare the test suite for a spin.
// Connect and prepare database for instance...
func (s *Suite) BeforeSuite(fn func()) {
	s.beforeSuiteHandlers = append(s.beforeSuiteHandlers, fn)
}

// BeforeFeature registers a function or method
// to be run once before every feature execution.
//
// If godog is run with concurrency option, it will
// run every feature per goroutine. So user may choose
// whether to isolate state within feature context or
// scenario.
//
// Best practice is not to have any state dependency on
// every scenario, but in some cases if VM for example
// needs to be started it may take very long for each
// scenario to restart it.
//
// Use it wisely and avoid sharing state between scenarios.
//
// Deprecated: BeforeFeature will be removed. Depending on
// your usecase, do setup in BeforeSuite or BeforeScenario.
func (s *Suite) BeforeFeature(fn func(*messages.GherkinDocument)) {
	s.beforeFeatureHandlers = append(s.beforeFeatureHandlers, fn)
}

// BeforeScenario registers a function or method
// to be run before every pickle.
//
// It is a good practice to restore the default state
// before every scenario so it would be isolated from
// any kind of state.
func (s *Suite) BeforeScenario(fn func(*messages.Pickle)) {
	s.beforeScenarioHandlers = append(s.beforeScenarioHandlers, fn)
}

// BeforeStep registers a function or method
// to be run before every step.
func (s *Suite) BeforeStep(fn func(*messages.Pickle_PickleStep)) {
	s.beforeStepHandlers = append(s.beforeStepHandlers, fn)
}

// AfterStep registers an function or method
// to be run after every step.
//
// It may be convenient to return a different kind of error
// in order to print more state details which may help
// in case of step failure
//
// In some cases, for example when running a headless
// browser, to take a screenshot after failure.
func (s *Suite) AfterStep(fn func(*messages.Pickle_PickleStep, error)) {
	s.afterStepHandlers = append(s.afterStepHandlers, fn)
}

// AfterScenario registers an function or method
// to be run after every pickle.
func (s *Suite) AfterScenario(fn func(*messages.Pickle, error)) {
	s.afterScenarioHandlers = append(s.afterScenarioHandlers, fn)
}

// AfterFeature registers a function or method
// to be run once after feature executed all scenarios.
//
// Deprecated: AfterFeature will be removed. Depending on
// your usecase, do cleanup and teardowns in AfterScenario
// or AfterSuite.
func (s *Suite) AfterFeature(fn func(*messages.GherkinDocument)) {
	s.afterFeatureHandlers = append(s.afterFeatureHandlers, fn)
}

// AfterSuite registers a function or method
// to be run once after suite runner
func (s *Suite) AfterSuite(fn func()) {
	s.afterSuiteHandlers = append(s.afterSuiteHandlers, fn)
}

func (s *Suite) run() {
	// run before suite handlers
	for _, f := range s.beforeSuiteHandlers {
		f()
	}
	// run features
	for _, f := range s.features {
		s.runFeature(f)
		if s.failed && s.stopOnFailure {
			// stop on first failure
			break
		}
	}
	// run after suite handlers
	for _, f := range s.afterSuiteHandlers {
		f()
	}
}

func (s *Suite) matchStep(step *messages.Pickle_PickleStep) *StepDefinition {
	def := s.matchStepText(step.Text)
	if def != nil && step.Argument != nil {
		def.args = append(def.args, step.Argument)
	}
	return def
}

func (s *Suite) runStep(pickle *messages.Pickle, step *messages.Pickle_PickleStep, prevStepErr error) (err error) {
	// run before step handlers
	for _, f := range s.beforeStepHandlers {
		f(step)
	}

	match := s.matchStep(step)
	s.fmt.Defined(pickle, step, match)

	// user multistep definitions may panic
	defer func() {
		if e := recover(); e != nil {
			err = &traceError{
				msg:   fmt.Sprintf("%v", e),
				stack: callStack(),
			}
		}

		if prevStepErr != nil {
			return
		}

		if err == ErrUndefined {
			return
		}

		switch err {
		case nil:
			s.fmt.Passed(pickle, step, match)
		case ErrPending:
			s.fmt.Pending(pickle, step, match)
		default:
			s.fmt.Failed(pickle, step, match, err)
		}

		// run after step handlers
		for _, f := range s.afterStepHandlers {
			f(step, err)
		}
	}()

	if undef, err := s.maybeUndefined(step.Text, step.Argument); err != nil {
		return err
	} else if len(undef) > 0 {
		if match != nil {
			match = &StepDefinition{
				args:      match.args,
				hv:        match.hv,
				Expr:      match.Expr,
				Handler:   match.Handler,
				nested:    match.nested,
				undefined: undef,
			}
		}
		s.fmt.Undefined(pickle, step, match)
		return ErrUndefined
	}

	if prevStepErr != nil {
		s.fmt.Skipped(pickle, step, match)
		return nil
	}

	err = s.maybeSubSteps(match.run())
	return
}

func (s *Suite) maybeUndefined(text string, arg interface{}) ([]string, error) {
	step := s.matchStepText(text)
	if nil == step {
		return []string{text}, nil
	}

	var undefined []string
	if !step.nested {
		return undefined, nil
	}

	if arg != nil {
		step.args = append(step.args, arg)
	}

	for _, next := range step.run().(Steps) {
		lines := strings.Split(next, "\n")
		// @TODO: we cannot currently parse table or content body from nested steps
		if len(lines) > 1 {
			return undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		if len(lines[0]) > 0 && lines[0][len(lines[0])-1] == ':' {
			return undefined, fmt.Errorf("nested steps cannot be multiline and have table or content body argument")
		}
		undef, err := s.maybeUndefined(next, nil)
		if err != nil {
			return undefined, err
		}
		undefined = append(undefined, undef...)
	}
	return undefined, nil
}

func (s *Suite) maybeSubSteps(result interface{}) error {
	if nil == result {
		return nil
	}

	if err, ok := result.(error); ok {
		return err
	}

	steps, ok := result.(Steps)
	if !ok {
		return fmt.Errorf("unexpected error, should have been []string: %T - %+v", result, result)
	}

	for _, text := range steps {
		if def := s.matchStepText(text); def == nil {
			return ErrUndefined
		} else if err := s.maybeSubSteps(def.run()); err != nil {
			return fmt.Errorf("%s: %+v", text, err)
		}
	}
	return nil
}

func (s *Suite) matchStepText(text string) *StepDefinition {
	for _, h := range s.steps {
		if m := h.Expr.FindStringSubmatch(text); len(m) > 0 {
			var args []interface{}
			for _, m := range m[1:] {
				args = append(args, m)
			}

			// since we need to assign arguments
			// better to copy the step definition
			return &StepDefinition{
				args:    args,
				hv:      h.hv,
				Expr:    h.Expr,
				Handler: h.Handler,
				nested:  h.nested,
			}
		}
	}
	return nil
}

func (s *Suite) runSteps(pickle *messages.Pickle, steps []*messages.Pickle_PickleStep) (err error) {
	for _, step := range steps {
		stepErr := s.runStep(pickle, step, err)
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
	return
}

func (s *Suite) shouldFail(err error) bool {
	if err == nil {
		return false
	}

	if err == ErrUndefined || err == ErrPending {
		return s.strict
	}

	return true
}

func (s *Suite) runFeature(f *feature) {
	if !isEmptyFeature(f.pickles) {
		for _, fn := range s.beforeFeatureHandlers {
			fn(f.GherkinDocument)
		}
	}

	s.fmt.Feature(f.GherkinDocument, f.Path, f.Content)

	defer func() {
		if !isEmptyFeature(f.pickles) {
			for _, fn := range s.afterFeatureHandlers {
				fn(f.GherkinDocument)
			}
		}
	}()

	for _, pickle := range f.pickles {
		err := s.runPickle(pickle)
		if s.shouldFail(err) {
			s.failed = true
			if s.stopOnFailure {
				return
			}
		}
	}
}

func isEmptyFeature(pickles []*messages.Pickle) bool {
	for _, pickle := range pickles {
		if len(pickle.Steps) > 0 {
			return false
		}
	}

	return true
}

func (s *Suite) runPickle(pickle *messages.Pickle) (err error) {
	if len(pickle.Steps) == 0 {
		s.fmt.Pickle(pickle)
		return ErrUndefined
	}

	// run before scenario handlers
	for _, f := range s.beforeScenarioHandlers {
		f(pickle)
	}

	s.fmt.Pickle(pickle)

	// scenario
	err = s.runSteps(pickle, pickle.Steps)

	// run after scenario handlers
	for _, f := range s.afterScenarioHandlers {
		f(pickle, err)
	}

	return
}

func (s *Suite) printStepDefinitions(w io.Writer) {
	var longest int
	for _, def := range s.steps {
		n := utf8.RuneCountInString(def.Expr.String())
		if longest < n {
			longest = n
		}
	}
	for _, def := range s.steps {
		n := utf8.RuneCountInString(def.Expr.String())
		location := def.definitionID()
		spaces := strings.Repeat(" ", longest-n)
		fmt.Fprintln(w, yellow(def.Expr.String())+spaces, blackb("# "+location))
	}
	if len(s.steps) == 0 {
		fmt.Fprintln(w, "there were no contexts registered, could not find any step definition..")
	}
}

var pathLineRe = regexp.MustCompile(`:([\d]+)$`)

func extractFeaturePathLine(p string) (string, int) {
	line := -1
	retPath := p
	if m := pathLineRe.FindStringSubmatch(p); len(m) > 0 {
		if i, err := strconv.Atoi(m[1]); err == nil {
			line = i
			retPath = p[:strings.LastIndexByte(p, ':')]
		}
	}
	return retPath, line
}

func parseFeatureFile(path string, newIDFunc func() string) (*feature, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	gherkinDocument, err := gherkin.ParseGherkinDocument(io.TeeReader(reader, &buf), newIDFunc)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", path, err)
	}

	pickles := gherkin.Pickles(*gherkinDocument, path, newIDFunc)

	return &feature{
		GherkinDocument: gherkinDocument,
		pickles:         pickles,
		Content:         buf.Bytes(),
		Path:            path,
	}, nil
}

func parseFeatureDir(dir string, newIDFunc func() string) ([]*feature, error) {
	var features []*feature
	return features, filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(p, ".feature") {
			return nil
		}

		feat, err := parseFeatureFile(p, newIDFunc)
		if err != nil {
			return err
		}
		features = append(features, feat)
		return nil
	})
}

func parsePath(path string) ([]*feature, error) {
	var features []*feature

	path, line := extractFeaturePathLine(path)

	fi, err := os.Stat(path)
	if err != nil {
		return features, err
	}

	newIDFunc := (&messages.Incrementing{}).NewId

	if fi.IsDir() {
		return parseFeatureDir(path, newIDFunc)
	}

	ft, err := parseFeatureFile(path, newIDFunc)
	if err != nil {
		return features, err
	}

	// filter scenario by line number
	var pickles []*messages.Pickle
	for _, pickle := range ft.pickles {
		sc := ft.findScenario(pickle.AstNodeIds[0])

		if line == -1 || uint32(line) == sc.Location.Line {
			pickles = append(pickles, pickle)
		}
	}
	ft.pickles = pickles

	return append(features, ft), nil
}

func parseFeatures(filter string, paths []string) ([]*feature, error) {
	byPath := make(map[string]*feature)
	var order int
	for _, path := range paths {
		feats, err := parsePath(path)
		switch {
		case os.IsNotExist(err):
			return nil, fmt.Errorf(`feature path "%s" is not available`, path)
		case os.IsPermission(err):
			return nil, fmt.Errorf(`feature path "%s" is not accessible`, path)
		case err != nil:
			return nil, err
		}

		for _, ft := range feats {
			if _, duplicate := byPath[ft.Path]; duplicate {
				continue
			}

			ft.order = order
			order++
			byPath[ft.Path] = ft
		}
	}

	return filterFeatures(filter, byPath), nil
}

type sortByOrderGiven []*feature

func (s sortByOrderGiven) Len() int           { return len(s) }
func (s sortByOrderGiven) Less(i, j int) bool { return s[i].order < s[j].order }
func (s sortByOrderGiven) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func filterFeatures(tags string, collected map[string]*feature) (features []*feature) {
	for _, ft := range collected {
		applyTagFilter(tags, ft)
		features = append(features, ft)
	}

	sort.Sort(sortByOrderGiven(features))

	return features
}

func applyTagFilter(tags string, ft *feature) {
	if len(tags) == 0 {
		return
	}

	var pickles []*messages.Pickle
	for _, pickle := range ft.pickles {
		if matchesTags(tags, pickle.Tags) {
			pickles = append(pickles, pickle)
		}
	}

	ft.pickles = pickles
}

// based on http://behat.readthedocs.org/en/v2.5/guides/6.cli.html#gherkin-filters
func matchesTags(filter string, tags []*messages.Pickle_PickleTag) (ok bool) {
	ok = true
	for _, andTags := range strings.Split(filter, "&&") {
		var okComma bool
		for _, tag := range strings.Split(andTags, ",") {
			tag = strings.Replace(strings.TrimSpace(tag), "@", "", -1)
			if tag[0] == '~' {
				tag = tag[1:]
				okComma = !hasTag(tags, tag) || okComma
			} else {
				okComma = hasTag(tags, tag) || okComma
			}
		}
		ok = ok && okComma
	}
	return
}

func hasTag(tags []*messages.Pickle_PickleTag, tag string) bool {
	for _, t := range tags {
		tName := strings.Replace(t.Name, "@", "", -1)

		if tName == tag {
			return true
		}
	}
	return false
}
