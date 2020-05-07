package godog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"

	"github.com/cucumber/godog/colors"

	"github.com/cucumber/messages-go/v10"
)

// some snippet formatting regexps
var snippetExprCleanup = regexp.MustCompile("([\\/\\[\\]\\(\\)\\\\^\\$\\.\\|\\?\\*\\+\\'])")
var snippetExprQuoted = regexp.MustCompile("(\\W|^)\"(?:[^\"]*)\"(\\W|$)")
var snippetMethodName = regexp.MustCompile("[^a-zA-Z\\_\\ ]")
var snippetNumbers = regexp.MustCompile("(\\d+)")

var snippetHelperFuncs = template.FuncMap{
	"backticked": func(s string) string {
		return "`" + s + "`"
	},
}

var undefinedSnippetsTpl = template.Must(template.New("snippets").Funcs(snippetHelperFuncs).Parse(`
{{ range . }}func {{ .Method }}({{ .Args }}) error {
	return godog.ErrPending
}

{{end}}func FeatureContext(s *godog.Suite) { {{ range . }}
	s.Step({{ backticked .Expr }}, {{ .Method }}){{end}}
}
`))

type undefinedSnippet struct {
	Method   string
	Expr     string
	argument *messages.PickleStepArgument
}

type registeredFormatter struct {
	name        string
	fmt         FormatterFunc
	description string
}

var formatters []*registeredFormatter

// FindFmt searches available formatters registered
// and returns FormaterFunc matched by given
// format name or nil otherwise
func FindFmt(name string) FormatterFunc {
	for _, el := range formatters {
		if el.name == name {
			return el.fmt
		}
	}
	return nil
}

// Format registers a feature suite output
// formatter by given name, description and
// FormatterFunc constructor function, to initialize
// formatter with the output recorder.
func Format(name, description string, f FormatterFunc) {
	formatters = append(formatters, &registeredFormatter{
		name:        name,
		fmt:         f,
		description: description,
	})
}

// AvailableFormatters gives a map of all
// formatters registered with their name as key
// and description as value
func AvailableFormatters() map[string]string {
	fmts := make(map[string]string, len(formatters))
	for _, f := range formatters {
		fmts[f.name] = f.description
	}
	return fmts
}

// Formatter is an interface for feature runner
// output summary presentation.
//
// New formatters may be created to represent
// suite results in different ways. These new
// formatters needs to be registered with a
// godog.Format function call
type Formatter interface {
	Feature(*messages.GherkinDocument, string, []byte)
	Pickle(*messages.Pickle)
	Defined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Failed(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition, error)
	Passed(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Skipped(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Undefined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Pending(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition)
	Summary()
}

// ConcurrentFormatter is an interface for a Concurrent
// version of the Formatter interface.
type ConcurrentFormatter interface {
	Formatter
	Copy(ConcurrentFormatter)
	Sync(ConcurrentFormatter)
}

// FormatterFunc builds a formatter with given
// suite name and io.Writer to record output
type FormatterFunc func(string, io.Writer) Formatter

type stepResultStatus int

const (
	passed stepResultStatus = iota
	failed
	skipped
	undefined
	pending
)

func (st stepResultStatus) clr() colors.ColorFunc {
	switch st {
	case passed:
		return green
	case failed:
		return red
	case skipped:
		return cyan
	default:
		return yellow
	}
}

func (st stepResultStatus) String() string {
	switch st {
	case passed:
		return "passed"
	case failed:
		return "failed"
	case skipped:
		return "skipped"
	case undefined:
		return "undefined"
	case pending:
		return "pending"
	default:
		return "unknown"
	}
}

type stepResult struct {
	status stepResultStatus
	time   time.Time
	err    error

	owner *messages.Pickle
	step  *messages.Pickle_PickleStep
	def   *StepDefinition
}

func newStepResult(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) *stepResult {
	return &stepResult{time: timeNowFunc(), owner: pickle, step: step, def: match}
}

func newBaseFmt(suite string, out io.Writer) *basefmt {
	return &basefmt{
		suiteName: suite,
		started:   timeNowFunc(),
		indent:    2,
		out:       out,
		lock:      new(sync.Mutex),
	}
}

type basefmt struct {
	suiteName string

	out    io.Writer
	owner  interface{}
	indent int

	started  time.Time
	features []*feature

	lock *sync.Mutex
}

func (f *basefmt) lastFeature() *feature {
	return f.features[len(f.features)-1]
}

func (f *basefmt) lastStepResult() *stepResult {
	return f.lastFeature().lastStepResult()
}

func (f *basefmt) findFeature(scenarioAstID string) *feature {
	for _, ft := range f.features {
		if sc := ft.findScenario(scenarioAstID); sc != nil {
			return ft
		}
	}

	panic("Couldn't find scenario for AST ID: " + scenarioAstID)
}

func (f *basefmt) findScenario(scenarioAstID string) *messages.GherkinDocument_Feature_Scenario {
	for _, ft := range f.features {
		if sc := ft.findScenario(scenarioAstID); sc != nil {
			return sc
		}
	}

	panic("Couldn't find scenario for AST ID: " + scenarioAstID)
}

func (f *basefmt) findBackground(scenarioAstID string) *messages.GherkinDocument_Feature_Background {
	for _, ft := range f.features {
		if bg := ft.findBackground(scenarioAstID); bg != nil {
			return bg
		}
	}

	return nil
}

func (f *basefmt) findExample(exampleAstID string) (*messages.GherkinDocument_Feature_Scenario_Examples, *messages.GherkinDocument_Feature_TableRow) {
	for _, ft := range f.features {
		if es, rs := ft.findExample(exampleAstID); es != nil && rs != nil {
			return es, rs
		}
	}

	return nil, nil
}

func (f *basefmt) findStep(stepAstID string) *messages.GherkinDocument_Feature_Step {
	for _, ft := range f.features {
		if st := ft.findStep(stepAstID); st != nil {
			return st
		}
	}

	panic("Couldn't find step for AST ID: " + stepAstID)
}

func (f *basefmt) Pickle(p *messages.Pickle) {
	f.lock.Lock()
	defer f.lock.Unlock()

	feature := f.features[len(f.features)-1]
	feature.pickleResults = append(feature.pickleResults, &pickleResult{Name: p.Name, time: timeNowFunc()})
}

func (f *basefmt) Defined(*messages.Pickle, *messages.Pickle_PickleStep, *StepDefinition) {}

func (f *basefmt) Feature(ft *messages.GherkinDocument, p string, c []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.features = append(f.features, &feature{Path: p, GherkinDocument: ft, time: timeNowFunc()})
}

func (f *basefmt) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.lock.Lock()
	defer f.lock.Unlock()

	s := newStepResult(pickle, step, match)
	s.status = passed
	f.lastFeature().appendStepResult(s)
}

func (f *basefmt) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.lock.Lock()
	defer f.lock.Unlock()

	s := newStepResult(pickle, step, match)
	s.status = skipped
	f.lastFeature().appendStepResult(s)
}

func (f *basefmt) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.lock.Lock()
	defer f.lock.Unlock()

	s := newStepResult(pickle, step, match)
	s.status = undefined
	f.lastFeature().appendStepResult(s)
}

func (f *basefmt) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	s := newStepResult(pickle, step, match)
	s.status = failed
	s.err = err
	f.lastFeature().appendStepResult(s)
}

func (f *basefmt) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.lock.Lock()
	defer f.lock.Unlock()

	s := newStepResult(pickle, step, match)
	s.status = pending
	f.lastFeature().appendStepResult(s)
}

func (f *basefmt) Summary() {
	var totalSc, passedSc, undefinedSc int
	var totalSt, passedSt, failedSt, skippedSt, pendingSt, undefinedSt int

	for _, feat := range f.features {
		for _, pr := range feat.pickleResults {
			var prStatus stepResultStatus
			totalSc++

			if len(pr.stepResults) == 0 {
				prStatus = undefined
			}

			for _, sr := range pr.stepResults {
				totalSt++

				switch sr.status {
				case passed:
					prStatus = passed
					passedSt++
				case failed:
					prStatus = failed
					failedSt++
				case skipped:
					skippedSt++
				case undefined:
					prStatus = undefined
					undefinedSt++
				case pending:
					prStatus = pending
					pendingSt++
				}
			}

			if prStatus == passed {
				passedSc++
			} else if prStatus == undefined {
				undefinedSc++
			}
		}
	}

	var steps, parts, scenarios []string
	if passedSt > 0 {
		steps = append(steps, green(fmt.Sprintf("%d passed", passedSt)))
	}
	if failedSt > 0 {
		parts = append(parts, red(fmt.Sprintf("%d failed", failedSt)))
		steps = append(steps, red(fmt.Sprintf("%d failed", failedSt)))
	}
	if pendingSt > 0 {
		parts = append(parts, yellow(fmt.Sprintf("%d pending", pendingSt)))
		steps = append(steps, yellow(fmt.Sprintf("%d pending", pendingSt)))
	}
	if undefinedSt > 0 {
		parts = append(parts, yellow(fmt.Sprintf("%d undefined", undefinedSc)))
		steps = append(steps, yellow(fmt.Sprintf("%d undefined", undefinedSt)))
	} else if undefinedSc > 0 {
		// there may be some scenarios without steps
		parts = append(parts, yellow(fmt.Sprintf("%d undefined", undefinedSc)))
	}
	if skippedSt > 0 {
		steps = append(steps, cyan(fmt.Sprintf("%d skipped", skippedSt)))
	}
	if passedSc > 0 {
		scenarios = append(scenarios, green(fmt.Sprintf("%d passed", passedSc)))
	}
	scenarios = append(scenarios, parts...)
	elapsed := timeNowFunc().Sub(f.started)

	fmt.Fprintln(f.out, "")

	if totalSc == 0 {
		fmt.Fprintln(f.out, "No scenarios")
	} else {
		fmt.Fprintln(f.out, fmt.Sprintf("%d scenarios (%s)", totalSc, strings.Join(scenarios, ", ")))
	}

	if totalSt == 0 {
		fmt.Fprintln(f.out, "No steps")
	} else {
		fmt.Fprintln(f.out, fmt.Sprintf("%d steps (%s)", totalSt, strings.Join(steps, ", ")))
	}

	elapsedString := elapsed.String()
	if elapsed.Nanoseconds() == 0 {
		// go 1.5 and 1.6 prints 0 instead of 0s, if duration is zero.
		elapsedString = "0s"
	}
	fmt.Fprintln(f.out, elapsedString)

	// prints used randomization seed
	seed, err := strconv.ParseInt(os.Getenv("GODOG_SEED"), 10, 64)
	if err == nil && seed != 0 {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, "Randomized with seed:", colors.Yellow(seed))
	}

	if text := f.snippets(); text != "" {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, yellow("You can implement step definitions for undefined steps with these snippets:"))
		fmt.Fprintln(f.out, yellow(text))
	}
}

func (f *basefmt) Sync(cf ConcurrentFormatter) {
	if source, ok := cf.(*basefmt); ok {
		f.lock = source.lock
	}
}

func (f *basefmt) Copy(cf ConcurrentFormatter) {
	if source, ok := cf.(*basefmt); ok {
		for _, v := range source.features {
			f.features = append(f.features, v)
		}
	}
}

func (s *undefinedSnippet) Args() (ret string) {
	var (
		args      []string
		pos       int
		breakLoop bool
	)
	for !breakLoop {
		part := s.Expr[pos:]
		ipos := strings.Index(part, "(\\d+)")
		spos := strings.Index(part, "\"([^\"]*)\"")
		switch {
		case spos == -1 && ipos == -1:
			breakLoop = true
		case spos == -1:
			pos += ipos + len("(\\d+)")
			args = append(args, reflect.Int.String())
		case ipos == -1:
			pos += spos + len("\"([^\"]*)\"")
			args = append(args, reflect.String.String())
		case ipos < spos:
			pos += ipos + len("(\\d+)")
			args = append(args, reflect.Int.String())
		case spos < ipos:
			pos += spos + len("\"([^\"]*)\"")
			args = append(args, reflect.String.String())
		}
	}

	if s.argument != nil {
		if s.argument.GetDocString() != nil {
			args = append(args, "*messages.PickleStepArgument_PickleDocString")
		}
		if s.argument.GetDataTable() != nil {
			args = append(args, "*messages.PickleStepArgument_PickleTable")
		}
	}

	var last string
	for i, arg := range args {
		if last == "" || last == arg {
			ret += fmt.Sprintf("arg%d, ", i+1)
		} else {
			ret = strings.TrimRight(ret, ", ") + fmt.Sprintf(" %s, arg%d, ", last, i+1)
		}
		last = arg
	}
	return strings.TrimSpace(strings.TrimRight(ret, ", ") + " " + last)
}

func (f *basefmt) findStepResults(status stepResultStatus) (res []*stepResult) {
	for _, feat := range f.features {
		for _, pr := range feat.pickleResults {
			for _, sr := range pr.stepResults {
				if sr.status == status {
					res = append(res, sr)
				}
			}
		}
	}

	return
}

func (f *basefmt) snippets() string {
	undefinedStepResults := f.findStepResults(undefined)
	if len(undefinedStepResults) == 0 {
		return ""
	}

	var index int
	var snips []*undefinedSnippet
	// build snippets
	for _, u := range undefinedStepResults {
		steps := []string{u.step.Text}
		arg := u.step.Argument
		if u.def != nil {
			steps = u.def.undefined
			arg = nil
		}
		for _, step := range steps {
			expr := snippetExprCleanup.ReplaceAllString(step, "\\$1")
			expr = snippetNumbers.ReplaceAllString(expr, "(\\d+)")
			expr = snippetExprQuoted.ReplaceAllString(expr, "$1\"([^\"]*)\"$2")
			expr = "^" + strings.TrimSpace(expr) + "$"

			name := snippetNumbers.ReplaceAllString(step, " ")
			name = snippetExprQuoted.ReplaceAllString(name, " ")
			name = strings.TrimSpace(snippetMethodName.ReplaceAllString(name, ""))
			var words []string
			for i, w := range strings.Split(name, " ") {
				switch {
				case i != 0:
					w = strings.Title(w)
				case len(w) > 0:
					w = string(unicode.ToLower(rune(w[0]))) + w[1:]
				}
				words = append(words, w)
			}
			name = strings.Join(words, "")
			if len(name) == 0 {
				index++
				name = fmt.Sprintf("StepDefinitioninition%d", index)
			}

			var found bool
			for _, snip := range snips {
				if snip.Expr == expr {
					found = true
					break
				}
			}
			if !found {
				snips = append(snips, &undefinedSnippet{Method: name, Expr: expr, argument: arg})
			}
		}
	}

	var buf bytes.Buffer
	if err := undefinedSnippetsTpl.Execute(&buf, snips); err != nil {
		panic(err)
	}
	// there may be trailing spaces
	return strings.Replace(buf.String(), " \n", "\n", -1)
}

func isLastStep(pickle *messages.Pickle, step *messages.Pickle_PickleStep) bool {
	return pickle.Steps[len(pickle.Steps)-1].Id == step.Id
}
