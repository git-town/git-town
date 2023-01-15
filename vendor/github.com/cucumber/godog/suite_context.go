package godog

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/gherkin-go/v11"
	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog/colors"
)

// SuiteContext provides steps for godog suite execution and
// can be used for meta-testing of godog features/steps themselves.
//
// Beware, steps or their definitions might change without backward
// compatibility guarantees. A typical user of the godog library should never
// need this, rather it is provided for those developing add-on libraries for godog.
//
// For an example of how to use, see godog's own `features/` and `suite_test.go`.
func SuiteContext(s *Suite, additionalContextInitializers ...func(suite *Suite)) {
	c := &suiteContext{
		extraCIs: additionalContextInitializers,
	}

	// apply any additional context intializers to modify the context that the
	// meta-tests will be run in
	for _, ci := range additionalContextInitializers {
		ci(s)
	}

	s.BeforeScenario(c.ResetBeforeEachScenario)

	s.Step(`^(?:a )?feature path "([^"]*)"$`, c.featurePath)
	s.Step(`^I parse features$`, c.parseFeatures)
	s.Step(`^I'm listening to suite events$`, c.iAmListeningToSuiteEvents)
	s.Step(`^I run feature suite$`, c.iRunFeatureSuite)
	s.Step(`^I run feature suite with tags "([^"]*)"$`, c.iRunFeatureSuiteWithTags)
	s.Step(`^I run feature suite with formatter "([^"]*)"$`, c.iRunFeatureSuiteWithFormatter)
	s.Step(`^(?:I )(allow|disable) variable injection`, c.iSetVariableInjectionTo)
	s.Step(`^(?:a )?feature "([^"]*)"(?: file)?:$`, c.aFeatureFile)
	s.Step(`^the suite should have (passed|failed)$`, c.theSuiteShouldHave)

	s.Step(`^I should have ([\d]+) features? files?:$`, c.iShouldHaveNumFeatureFiles)
	s.Step(`^I should have ([\d]+) scenarios? registered$`, c.numScenariosRegistered)
	s.Step(`^there (was|were) ([\d]+) "([^"]*)" events? fired$`, c.thereWereNumEventsFired)
	s.Step(`^there was event triggered before scenario "([^"]*)"$`, c.thereWasEventTriggeredBeforeScenario)
	s.Step(`^these events had to be fired for a number of times:$`, c.theseEventsHadToBeFiredForNumberOfTimes)

	s.Step(`^(?:a )?failing step`, c.aFailingStep)
	s.Step(`^this step should fail`, c.aFailingStep)
	s.Step(`^the following steps? should be (passed|failed|skipped|undefined|pending):`, c.followingStepsShouldHave)
	s.Step(`^all steps should (?:be|have|have been) (passed|failed|skipped|undefined|pending)$`, c.allStepsShouldHave)
	s.Step(`^the undefined step snippets should be:$`, c.theUndefinedStepSnippetsShouldBe)

	// event stream
	s.Step(`^the following events should be fired:$`, c.thereShouldBeEventsFired)

	// lt
	s.Step(`^savybių aplankas "([^"]*)"$`, c.featurePath)
	s.Step(`^aš išskaitau savybes$`, c.parseFeatures)
	s.Step(`^aš turėčiau turėti ([\d]+) savybių failus:$`, c.iShouldHaveNumFeatureFiles)

	s.Step(`^(?:a )?pending step$`, func() error {
		return ErrPending
	})
	s.Step(`^(?:a )?passing step$`, func() error {
		return nil
	})

	// Introduced to test formatter/cucumber.feature
	s.Step(`^the rendered json will be as follows:$`, c.theRenderJSONWillBe)

	// Introduced to test formatter/pretty.feature
	s.Step(`^the rendered output will be as follows:$`, c.theRenderOutputWillBe)

	// Introduced to test formatter/junit.feature
	s.Step(`^the rendered xml will be as follows:$`, c.theRenderXMLWillBe)

	s.Step(`^(?:a )?failing multistep$`, func() Steps {
		return Steps{"passing step", "failing step"}
	})

	s.Step(`^(?:a |an )?undefined multistep$`, func() Steps {
		return Steps{"passing step", "undefined step", "passing step"}
	})

	s.Step(`^(?:a )?passing multistep$`, func() Steps {
		return Steps{"passing step", "passing step", "passing step"}
	})

	s.Step(`^(?:a )?failing nested multistep$`, func() Steps {
		return Steps{"passing step", "passing multistep", "failing multistep"}
	})
	// Default recovery step
	s.Step(`Ignore.*`, func() error {
		return nil
	})

	s.BeforeStep(c.inject)
}

func (s *suiteContext) inject(step *messages.Pickle_PickleStep) {
	if !s.allowInjection {
		return
	}

	step.Text = injectAll(step.Text)

	if table := step.Argument.GetDataTable(); table != nil {
		for i := 0; i < len(table.Rows); i++ {
			for n, cell := range table.Rows[i].Cells {
				table.Rows[i].Cells[n].Value = injectAll(cell.Value)
			}
		}
	}

	if doc := step.Argument.GetDocString(); doc != nil {
		doc.Content = injectAll(doc.Content)
	}
}

func injectAll(src string) string {
	re := regexp.MustCompile(`{{[^{}]+}}`)
	return re.ReplaceAllStringFunc(
		src,
		func(key string) string {
			injectRegex := regexp.MustCompile(`^{{.+}}$`)
			if injectRegex.MatchString(key) {
				return "someverylonginjectionsoweacanbesureitsurpasstheinitiallongeststeplenghtanditwillhelptestsmethodsafety"
			}
			return key
		},
	)
}

type firedEvent struct {
	name string
	args []interface{}
}

type suiteContext struct {
	paths          []string
	testedSuite    *Suite
	extraCIs       []func(suite *Suite)
	events         []*firedEvent
	out            bytes.Buffer
	allowInjection bool
}

func (s *suiteContext) ResetBeforeEachScenario(*messages.Pickle) {
	// reset whole suite with the state
	s.out.Reset()
	s.paths = []string{}
	s.testedSuite = &Suite{}
	// our tested suite will have the same context registered
	SuiteContext(s.testedSuite, s.extraCIs...)
	// reset all fired events
	s.events = []*firedEvent{}
	s.allowInjection = false
}

func (s *suiteContext) iSetVariableInjectionTo(to string) error {
	s.allowInjection = to == "allow"
	return nil
}

func (s *suiteContext) iRunFeatureSuiteWithTags(tags string) error {
	if err := s.parseFeatures(); err != nil {
		return err
	}
	for _, feat := range s.testedSuite.features {
		applyTagFilter(tags, feat)
	}
	s.testedSuite.fmt = testFormatterFunc("godog", &s.out)
	s.testedSuite.run()
	s.testedSuite.fmt.Summary()
	return nil
}

func (s *suiteContext) iRunFeatureSuiteWithFormatter(name string) error {
	f := FindFmt(name)
	if f == nil {
		return fmt.Errorf(`formatter "%s" is not available`, name)
	}
	s.testedSuite.fmt = f("godog", colors.Uncolored(&s.out))
	if err := s.parseFeatures(); err != nil {
		return err
	}
	s.testedSuite.run()
	s.testedSuite.fmt.Summary()
	return nil
}

func (s *suiteContext) thereShouldBeEventsFired(doc *messages.PickleStepArgument_PickleDocString) error {
	actual := strings.Split(strings.TrimSpace(s.out.String()), "\n")
	expect := strings.Split(strings.TrimSpace(doc.Content), "\n")
	if len(expect) != len(actual) {
		return fmt.Errorf("expected %d events, but got %d", len(expect), len(actual))
	}

	type ev struct {
		Event string
	}

	for i, event := range actual {
		exp := strings.TrimSpace(expect[i])
		var act ev
		if err := json.Unmarshal([]byte(event), &act); err != nil {
			return fmt.Errorf("failed to read event data: %v", err)
		}

		if act.Event != exp {
			return fmt.Errorf(`expected event: "%s" at position: %d, but actual was "%s"`, exp, i, act.Event)
		}
	}
	return nil
}

func (s *suiteContext) cleanupSnippet(snip string) string {
	lines := strings.Split(strings.TrimSpace(snip), "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func (s *suiteContext) theUndefinedStepSnippetsShouldBe(body *messages.PickleStepArgument_PickleDocString) error {
	f, ok := s.testedSuite.fmt.(*testFormatter)
	if !ok {
		return fmt.Errorf("this step requires testFormatter, but there is: %T", s.testedSuite.fmt)
	}
	actual := s.cleanupSnippet(f.snippets())
	expected := s.cleanupSnippet(body.Content)
	if actual != expected {
		return fmt.Errorf("snippets do not match actual: %s", f.snippets())
	}
	return nil
}

func (s *suiteContext) followingStepsShouldHave(status string, steps *messages.PickleStepArgument_PickleDocString) error {
	var expected = strings.Split(steps.Content, "\n")
	var actual, unmatched, matched []string

	f, ok := s.testedSuite.fmt.(*testFormatter)
	if !ok {
		return fmt.Errorf("this step requires testFormatter, but there is: %T", s.testedSuite.fmt)
	}
	switch status {
	case "passed":
		for _, st := range f.findStepResults(passed) {
			actual = append(actual, st.step.Text)
		}
	case "failed":
		for _, st := range f.findStepResults(failed) {
			actual = append(actual, st.step.Text)
		}
	case "skipped":
		for _, st := range f.findStepResults(skipped) {
			actual = append(actual, st.step.Text)
		}
	case "undefined":
		for _, st := range f.findStepResults(undefined) {
			actual = append(actual, st.step.Text)
		}
	case "pending":
		for _, st := range f.findStepResults(pending) {
			actual = append(actual, st.step.Text)
		}
	default:
		return fmt.Errorf("unexpected step status wanted: %s", status)
	}

	if len(expected) > len(actual) {
		return fmt.Errorf("number of expected %s steps: %d is less than actual %s steps: %d", status, len(expected), status, len(actual))
	}

	for _, a := range actual {
		for _, e := range expected {
			if a == e {
				matched = append(matched, e)
				break
			}
		}
	}

	if len(matched) >= len(expected) {
		return nil
	}
	for _, s := range expected {
		var found bool
		for _, m := range matched {
			if s == m {
				found = true
				break
			}
		}
		if !found {
			unmatched = append(unmatched, s)
		}
	}

	return fmt.Errorf("the steps: %s - are not %s", strings.Join(unmatched, ", "), status)
}

func (s *suiteContext) allStepsShouldHave(status string) error {
	f, ok := s.testedSuite.fmt.(*testFormatter)
	if !ok {
		return fmt.Errorf("this step requires testFormatter, but there is: %T", s.testedSuite.fmt)
	}

	total := len(f.findStepResults(passed)) +
		len(f.findStepResults(failed)) +
		len(f.findStepResults(skipped)) +
		len(f.findStepResults(undefined)) +
		len(f.findStepResults(pending))
	var actual int
	switch status {
	case "passed":
		actual = len(f.findStepResults(passed))
	case "failed":
		actual = len(f.findStepResults(failed))
	case "skipped":
		actual = len(f.findStepResults(skipped))
	case "undefined":
		actual = len(f.findStepResults(undefined))
	case "pending":
		actual = len(f.findStepResults(pending))
	default:
		return fmt.Errorf("unexpected step status wanted: %s", status)
	}

	if total > actual {
		return fmt.Errorf("number of expected %s steps: %d is less than actual %s steps: %d", status, total, status, actual)
	}
	return nil
}

func (s *suiteContext) iAmListeningToSuiteEvents() error {
	s.testedSuite.BeforeSuite(func() {
		s.events = append(s.events, &firedEvent{"BeforeSuite", []interface{}{}})
	})
	s.testedSuite.AfterSuite(func() {
		s.events = append(s.events, &firedEvent{"AfterSuite", []interface{}{}})
	})
	s.testedSuite.BeforeFeature(func(ft *messages.GherkinDocument) {
		s.events = append(s.events, &firedEvent{"BeforeFeature", []interface{}{ft}})
	})
	s.testedSuite.AfterFeature(func(ft *messages.GherkinDocument) {
		s.events = append(s.events, &firedEvent{"AfterFeature", []interface{}{ft}})
	})
	s.testedSuite.BeforeScenario(func(pickle *messages.Pickle) {
		s.events = append(s.events, &firedEvent{"BeforeScenario", []interface{}{pickle}})
	})
	s.testedSuite.AfterScenario(func(pickle *messages.Pickle, err error) {
		s.events = append(s.events, &firedEvent{"AfterScenario", []interface{}{pickle, err}})
	})
	s.testedSuite.BeforeStep(func(step *messages.Pickle_PickleStep) {
		s.events = append(s.events, &firedEvent{"BeforeStep", []interface{}{step}})
	})
	s.testedSuite.AfterStep(func(step *messages.Pickle_PickleStep, err error) {
		s.events = append(s.events, &firedEvent{"AfterStep", []interface{}{step, err}})
	})
	return nil
}

func (s *suiteContext) aFailingStep() error {
	return fmt.Errorf("intentional failure")
}

// parse a given feature file body as a feature
func (s *suiteContext) aFeatureFile(path string, body *messages.PickleStepArgument_PickleDocString) error {
	gd, err := gherkin.ParseGherkinDocument(strings.NewReader(body.Content), (&messages.Incrementing{}).NewId)
	pickles := gherkin.Pickles(*gd, path, (&messages.Incrementing{}).NewId)
	s.testedSuite.features = append(s.testedSuite.features, &feature{GherkinDocument: gd, pickles: pickles, Path: path})
	return err
}

func (s *suiteContext) featurePath(path string) error {
	s.paths = append(s.paths, path)
	return nil
}

func (s *suiteContext) parseFeatures() error {
	fts, err := parseFeatures("", s.paths)
	if err != nil {
		return err
	}
	s.testedSuite.features = append(s.testedSuite.features, fts...)
	return nil
}

func (s *suiteContext) theSuiteShouldHave(state string) error {
	if s.testedSuite.failed && state == "passed" {
		return fmt.Errorf("the feature suite has failed")
	}
	if !s.testedSuite.failed && state == "failed" {
		return fmt.Errorf("the feature suite has passed")
	}
	return nil
}

func (s *suiteContext) iShouldHaveNumFeatureFiles(num int, files *messages.PickleStepArgument_PickleDocString) error {
	if len(s.testedSuite.features) != num {
		return fmt.Errorf("expected %d features to be parsed, but have %d", num, len(s.testedSuite.features))
	}
	expected := strings.Split(files.Content, "\n")
	var actual []string
	for _, ft := range s.testedSuite.features {
		actual = append(actual, ft.Path)
	}
	if len(expected) != len(actual) {
		return fmt.Errorf("expected %d feature paths to be parsed, but have %d", len(expected), len(actual))
	}
	for i := 0; i < len(expected); i++ {
		var matched bool
		split := strings.Split(expected[i], "/")
		exp := filepath.Join(split...)
		for j := 0; j < len(actual); j++ {
			split = strings.Split(actual[j], "/")
			act := filepath.Join(split...)
			if exp == act {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf(`expected feature path "%s" at position: %d, was not parsed, actual are %+v`, exp, i, actual)
		}
	}
	return nil
}

func (s *suiteContext) iRunFeatureSuite() error {
	if err := s.parseFeatures(); err != nil {
		return err
	}
	s.testedSuite.fmt = testFormatterFunc("godog", &s.out)
	s.testedSuite.run()
	s.testedSuite.fmt.Summary()

	return nil
}

func (s *suiteContext) numScenariosRegistered(expected int) (err error) {
	var num int
	for _, ft := range s.testedSuite.features {
		num += len(ft.pickles)
	}
	if num != expected {
		err = fmt.Errorf("expected %d scenarios to be registered, but got %d", expected, num)
	}
	return
}

func (s *suiteContext) thereWereNumEventsFired(_ string, expected int, typ string) error {
	var num int
	for _, event := range s.events {
		if event.name == typ {
			num++
		}
	}
	if num != expected {
		return fmt.Errorf("expected %d %s events to be fired, but got %d", expected, typ, num)
	}
	return nil
}

func (s *suiteContext) thereWasEventTriggeredBeforeScenario(expected string) error {
	var found []string
	for _, event := range s.events {
		if event.name != "BeforeScenario" {
			continue
		}

		var name string
		switch t := event.args[0].(type) {
		case *messages.Pickle:
			name = t.Name
		}
		if name == expected {
			return nil
		}

		found = append(found, name)
	}

	if len(found) == 0 {
		return fmt.Errorf("before scenario event was never triggered or listened")
	}

	return fmt.Errorf(`expected "%s" scenario, but got these fired %s`, expected, `"`+strings.Join(found, `", "`)+`"`)
}

func (s *suiteContext) theseEventsHadToBeFiredForNumberOfTimes(tbl *messages.PickleStepArgument_PickleTable) error {
	if len(tbl.Rows[0].Cells) != 2 {
		return fmt.Errorf("expected two columns for event table row, got: %d", len(tbl.Rows[0].Cells))
	}

	for _, row := range tbl.Rows {
		num, err := strconv.ParseInt(row.Cells[1].Value, 10, 0)
		if err != nil {
			return err
		}
		if err := s.thereWereNumEventsFired("", int(num), row.Cells[0].Value); err != nil {
			return err
		}
	}
	return nil
}

func (s *suiteContext) theRenderJSONWillBe(docstring *messages.PickleStepArgument_PickleDocString) error {
	suiteCtxReg := regexp.MustCompile(`suite_context.go:\d+`)

	expectedString := docstring.Content
	expectedString = suiteCtxReg.ReplaceAllString(expectedString, `suite_context.go:0`)
	actualString := s.out.String()
	actualString = suiteCtxReg.ReplaceAllString(actualString, `suite_context.go:0`)

	var expected []cukeFeatureJSON
	if err := json.Unmarshal([]byte(expectedString), &expected); err != nil {
		return err
	}

	var actual []cukeFeatureJSON
	if err := json.Unmarshal([]byte(actualString), &actual); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected json does not match actual: %s", actualString)
	}
	return nil
}

func (s *suiteContext) theRenderOutputWillBe(docstring *messages.PickleStepArgument_PickleDocString) error {
	suiteCtxReg := regexp.MustCompile(`suite_context.go:\d+`)
	suiteCtxFuncReg := regexp.MustCompile(`github.com/cucumber/godog.SuiteContext.func(\d+)`)

	expected := trimAllLines(strings.TrimSpace(docstring.Content))
	expected = suiteCtxReg.ReplaceAllString(expected, "suite_context.go:0")
	expected = suiteCtxFuncReg.ReplaceAllString(expected, "SuiteContext.func$1")

	actual := trimAllLines(strings.TrimSpace(s.out.String()))
	actual = suiteCtxReg.ReplaceAllString(actual, "suite_context.go:0")
	actual = suiteCtxFuncReg.ReplaceAllString(actual, "SuiteContext.func$1")

	if err := match(expected, actual); err != nil {
		return err
	}

	return nil
}

func (s *suiteContext) theRenderXMLWillBe(docstring *messages.PickleStepArgument_PickleDocString) error {
	expectedString := docstring.Content
	actualString := s.out.String()

	var expected junitPackageSuite
	if err := xml.Unmarshal([]byte(expectedString), &expected); err != nil {
		return err
	}

	var actual junitPackageSuite
	if err := xml.Unmarshal([]byte(actualString), &actual); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected json does not match actual: %s", actualString)
	}
	return nil
}

type testFormatter struct {
	*basefmt
	pickles []*messages.Pickle
}

func testFormatterFunc(suite string, out io.Writer) Formatter {
	return &testFormatter{basefmt: newBaseFmt(suite, out)}
}

func (f *testFormatter) Pickle(p *messages.Pickle) {
	f.basefmt.Pickle(p)
	f.pickles = append(f.pickles, p)
}

func (f *testFormatter) Summary() {}

func match(expected, actual string) error {
	act := []byte(actual)
	exp := []byte(expected)

	if len(act) != len(exp) {
		return fmt.Errorf("content lengths do not match, expected: %d, actual %d, expected output:\n%s, actual output:\n%s", len(exp), len(act), expected, actual)
	}

	for i := 0; i < len(exp); i++ {
		if act[i] == exp[i] {
			continue
		}

		cpe := make([]byte, len(exp))
		copy(cpe, exp)
		e := append(exp[:i], '^')
		e = append(e, cpe[i:]...)

		cpa := make([]byte, len(act))
		copy(cpa, act)
		a := append(act[:i], '^')
		a = append(a, cpa[i:]...)

		return fmt.Errorf("expected output does not match:\n%s\n\n%s", string(a), string(e))
	}

	return nil
}
