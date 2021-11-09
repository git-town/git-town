package formatters

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/colors"
	"github.com/cucumber/godog/formatters"
)

func init() {
	formatters.Format("pretty", "Prints every feature with runtime statuses.", PrettyFormatterFunc)
}

// PrettyFormatterFunc implements the FormatterFunc for the pretty formatter
func PrettyFormatterFunc(suite string, out io.Writer) formatters.Formatter {
	return &Pretty{Base: NewBase(suite, out)}
}

var outlinePlaceholderRegexp = regexp.MustCompile("<[^>]+>")

// Pretty is a formatter for readable output.
type Pretty struct {
	*Base
	firstFeature *bool
}

// TestRunStarted is triggered on test start.
func (f *Pretty) TestRunStarted() {
	f.Base.TestRunStarted()

	f.Lock.Lock()
	defer f.Lock.Unlock()

	firstFeature := true
	f.firstFeature = &firstFeature
}

// Feature receives gherkin document.
func (f *Pretty) Feature(gd *messages.GherkinDocument, p string, c []byte) {
	f.Lock.Lock()
	if !*f.firstFeature {
		fmt.Fprintln(f.out, "")
	}

	*f.firstFeature = false
	f.Lock.Unlock()

	f.Base.Feature(gd, p, c)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printFeature(gd.Feature)
}

// Pickle takes a gherkin node for formatting.
func (f *Pretty) Pickle(pickle *messages.Pickle) {
	f.Base.Pickle(pickle)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	if len(pickle.Steps) == 0 {
		f.printUndefinedPickle(pickle)
		return
	}
}

// Passed captures passed step.
func (f *Pretty) Passed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Passed(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printStep(pickle, step)
}

// Skipped captures skipped step.
func (f *Pretty) Skipped(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Skipped(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printStep(pickle, step)
}

// Undefined captures undefined step.
func (f *Pretty) Undefined(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Undefined(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printStep(pickle, step)
}

// Failed captures failed step.
func (f *Pretty) Failed(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition, err error) {
	f.Base.Failed(pickle, step, match, err)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printStep(pickle, step)
}

// Pending captures pending step.
func (f *Pretty) Pending(pickle *messages.Pickle, step *messages.PickleStep, match *formatters.StepDefinition) {
	f.Base.Pending(pickle, step, match)

	f.Lock.Lock()
	defer f.Lock.Unlock()

	f.printStep(pickle, step)
}

func (f *Pretty) printFeature(feature *messages.Feature) {
	fmt.Fprintln(f.out, keywordAndName(feature.Keyword, feature.Name))
	if strings.TrimSpace(feature.Description) != "" {
		for _, line := range strings.Split(feature.Description, "\n") {
			fmt.Fprintln(f.out, s(f.indent)+strings.TrimSpace(line))
		}
	}
}

func keywordAndName(keyword, name string) string {
	title := whiteb(keyword + ":")
	if len(name) > 0 {
		title += " " + name
	}
	return title
}

func (f *Pretty) scenarioLengths(pickle *messages.Pickle) (scenarioHeaderLength int, maxLength int) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	astScenario := feature.FindScenario(pickle.AstNodeIds[0])
	astBackground := feature.FindBackground(pickle.AstNodeIds[0])

	scenarioHeaderLength = f.lengthPickle(astScenario.Keyword, astScenario.Name)
	maxLength = f.longestStep(astScenario.Steps, scenarioHeaderLength)

	if astBackground != nil {
		maxLength = f.longestStep(astBackground.Steps, maxLength)
	}

	return scenarioHeaderLength, maxLength
}

func (f *Pretty) printScenarioHeader(pickle *messages.Pickle, astScenario *messages.Scenario, spaceFilling int) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	text := s(f.indent) + keywordAndName(astScenario.Keyword, astScenario.Name)
	text += s(spaceFilling) + line(feature.Uri, astScenario.Location)
	fmt.Fprintln(f.out, "\n"+text)
}

func (f *Pretty) printUndefinedPickle(pickle *messages.Pickle) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	astScenario := feature.FindScenario(pickle.AstNodeIds[0])
	astBackground := feature.FindBackground(pickle.AstNodeIds[0])

	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)

	if astBackground != nil {
		fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
		for _, step := range astBackground.Steps {
			text := s(f.indent*2) + cyan(strings.TrimSpace(step.Keyword)) + " " + cyan(step.Text)
			fmt.Fprintln(f.out, text)
		}
	}

	//  do not print scenario headers and examples multiple times
	if len(astScenario.Examples) > 0 {
		exampleTable, exampleRow := feature.FindExample(pickle.AstNodeIds[1])
		firstExampleRow := exampleTable.TableBody[0].Id == exampleRow.Id
		firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

		if !(firstExamplesTable && firstExampleRow) {
			return
		}
	}

	f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)

	for _, examples := range astScenario.Examples {
		max := longestExampleRow(examples, cyan, cyan)

		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, s(f.indent*2)+keywordAndName(examples.Keyword, examples.Name))

		f.printTableHeader(examples.TableHeader, max)

		for _, row := range examples.TableBody {
			f.printTableRow(row, max, cyan)
		}
	}
}

// Summary renders summary information.
func (f *Pretty) Summary() {
	failedStepResults := f.Storage.MustGetPickleStepResultsByStatus(failed)
	if len(failedStepResults) > 0 {
		fmt.Fprintln(f.out, "\n--- "+red("Failed steps:")+"\n")

		sort.Sort(sortPickleStepResultsByPickleStepID(failedStepResults))

		for _, fail := range failedStepResults {
			pickle := f.Storage.MustGetPickle(fail.PickleID)
			pickleStep := f.Storage.MustGetPickleStep(fail.PickleStepID)
			feature := f.Storage.MustGetFeature(pickle.Uri)

			astScenario := feature.FindScenario(pickle.AstNodeIds[0])
			scenarioDesc := fmt.Sprintf("%s: %s", astScenario.Keyword, pickle.Name)

			astStep := feature.FindStep(pickleStep.AstNodeIds[0])
			stepDesc := strings.TrimSpace(astStep.Keyword) + " " + pickleStep.Text

			fmt.Fprintln(f.out, s(f.indent)+red(scenarioDesc)+line(feature.Uri, astScenario.Location))
			fmt.Fprintln(f.out, s(f.indent*2)+red(stepDesc)+line(feature.Uri, astStep.Location))
			fmt.Fprintln(f.out, s(f.indent*3)+red("Error: ")+redb(fmt.Sprintf("%+v", fail.Err))+"\n")
		}
	}

	f.Base.Summary()
}

func (f *Pretty) printOutlineExample(pickle *messages.Pickle, backgroundSteps int) {
	var errorMsg string
	var clr = green

	feature := f.Storage.MustGetFeature(pickle.Uri)
	astScenario := feature.FindScenario(pickle.AstNodeIds[0])
	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)

	exampleTable, exampleRow := feature.FindExample(pickle.AstNodeIds[1])
	printExampleHeader := exampleTable.TableBody[0].Id == exampleRow.Id
	firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

	pickleStepResults := f.Storage.MustGetPickleStepResultsByPickleID(pickle.Id)

	firstExecutedScenarioStep := len(pickleStepResults) == backgroundSteps+1
	if firstExamplesTable && printExampleHeader && firstExecutedScenarioStep {
		f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)
	}

	if len(exampleTable.TableBody) == 0 {
		// do not print empty examples
		return
	}

	lastStep := len(pickleStepResults) == len(pickle.Steps)
	if !lastStep {
		// do not print examples unless all steps has finished
		return
	}

	for _, result := range pickleStepResults {
		// determine example row status
		switch {
		case result.Status == failed:
			errorMsg = result.Err.Error()
			clr = result.Status.Color()
		case result.Status == undefined || result.Status == pending:
			clr = result.Status.Color()
		case result.Status == skipped && clr == nil:
			clr = cyan
		}

		if firstExamplesTable && printExampleHeader {
			// in first example, we need to print steps

			pickleStep := f.Storage.MustGetPickleStep(result.PickleStepID)
			astStep := feature.FindStep(pickleStep.AstNodeIds[0])

			var text = ""
			if result.Def != nil {
				if m := outlinePlaceholderRegexp.FindAllStringIndex(astStep.Text, -1); len(m) > 0 {
					var pos int
					for i := 0; i < len(m); i++ {
						pair := m[i]
						text += cyan(astStep.Text[pos:pair[0]])
						text += cyanb(astStep.Text[pair[0]:pair[1]])
						pos = pair[1]
					}
					text += cyan(astStep.Text[pos:len(astStep.Text)])
				} else {
					text = cyan(astStep.Text)
				}

				_, maxLength := f.scenarioLengths(pickle)
				stepLength := f.lengthPickleStep(astStep.Keyword, astStep.Text)

				text += s(maxLength - stepLength)
				text += " " + blackb("# "+DefinitionID(result.Def))
			}

			// print the step outline
			fmt.Fprintln(f.out, s(f.indent*2)+cyan(strings.TrimSpace(astStep.Keyword))+" "+text)

			if pickleStep.Argument != nil {
				if table := pickleStep.Argument.DataTable; table != nil {
					f.printTable(table, cyan)
				}

				if docString := astStep.DocString; docString != nil {
					f.printDocString(docString)
				}
			}
		}
	}

	max := longestExampleRow(exampleTable, clr, cyan)

	// an example table header
	if printExampleHeader {
		fmt.Fprintln(f.out, "")
		fmt.Fprintln(f.out, s(f.indent*2)+keywordAndName(exampleTable.Keyword, exampleTable.Name))

		f.printTableHeader(exampleTable.TableHeader, max)
	}

	f.printTableRow(exampleRow, max, clr)

	if errorMsg != "" {
		fmt.Fprintln(f.out, s(f.indent*4)+redb(errorMsg))
	}
}

func (f *Pretty) printTableRow(row *messages.TableRow, max []int, clr colors.ColorFunc) {
	cells := make([]string, len(row.Cells))

	for i, cell := range row.Cells {
		val := clr(cell.Value)
		ln := utf8.RuneCountInString(val)
		cells[i] = val + s(max[i]-ln)
	}

	fmt.Fprintln(f.out, s(f.indent*3)+"| "+strings.Join(cells, " | ")+" |")
}

func (f *Pretty) printTableHeader(row *messages.TableRow, max []int) {
	f.printTableRow(row, max, cyan)
}

func (f *Pretty) printStep(pickle *messages.Pickle, pickleStep *messages.PickleStep) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	astBackground := feature.FindBackground(pickle.AstNodeIds[0])
	astScenario := feature.FindScenario(pickle.AstNodeIds[0])
	astStep := feature.FindStep(pickleStep.AstNodeIds[0])

	var astBackgroundStep bool
	var firstExecutedBackgroundStep bool
	var backgroundSteps int
	if astBackground != nil {
		backgroundSteps = len(astBackground.Steps)

		for idx, step := range astBackground.Steps {
			if step.Id == pickleStep.AstNodeIds[0] {
				astBackgroundStep = true
				firstExecutedBackgroundStep = idx == 0
				break
			}
		}
	}

	firstPickle := feature.Pickles[0].Id == pickle.Id

	if astBackgroundStep && !firstPickle {
		return
	}

	if astBackgroundStep && firstExecutedBackgroundStep {
		fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
	}

	if !astBackgroundStep && len(astScenario.Examples) > 0 {
		f.printOutlineExample(pickle, backgroundSteps)
		return
	}

	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle)
	stepLength := f.lengthPickleStep(astStep.Keyword, pickleStep.Text)

	firstExecutedScenarioStep := astScenario.Steps[0].Id == pickleStep.AstNodeIds[0]
	if !astBackgroundStep && firstExecutedScenarioStep {
		f.printScenarioHeader(pickle, astScenario, maxLength-scenarioHeaderLength)
	}

	pickleStepResult := f.Storage.MustGetPickleStepResult(pickleStep.Id)
	text := s(f.indent*2) + pickleStepResult.Status.Color()(strings.TrimSpace(astStep.Keyword)) + " " + pickleStepResult.Status.Color()(pickleStep.Text)
	if pickleStepResult.Def != nil {
		text += s(maxLength - stepLength + 1)
		text += blackb("# " + DefinitionID(pickleStepResult.Def))
	}
	fmt.Fprintln(f.out, text)

	if pickleStep.Argument != nil {
		if table := pickleStep.Argument.DataTable; table != nil {
			f.printTable(table, cyan)
		}

		if docString := astStep.DocString; docString != nil {
			f.printDocString(docString)
		}
	}

	if pickleStepResult.Err != nil {
		fmt.Fprintln(f.out, s(f.indent*2)+redb(fmt.Sprintf("%+v", pickleStepResult.Err)))
	}

	if pickleStepResult.Status == pending {
		fmt.Fprintln(f.out, s(f.indent*3)+yellow("TODO: write pending definition"))
	}
}

func (f *Pretty) printDocString(docString *messages.DocString) {
	var ct string

	if len(docString.MediaType) > 0 {
		ct = " " + cyan(docString.MediaType)
	}

	fmt.Fprintln(f.out, s(f.indent*3)+cyan(docString.Delimiter)+ct)

	for _, ln := range strings.Split(docString.Content, "\n") {
		fmt.Fprintln(f.out, s(f.indent*3)+cyan(ln))
	}

	fmt.Fprintln(f.out, s(f.indent*3)+cyan(docString.Delimiter))
}

// print table with aligned table cells
// @TODO: need to make example header cells bold
func (f *Pretty) printTable(t *messages.PickleTable, c colors.ColorFunc) {
	maxColLengths := maxColLengths(t, c)
	var cols = make([]string, len(t.Rows[0].Cells))

	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			val := c(cell.Value)
			colLength := utf8.RuneCountInString(val)
			cols[i] = val + s(maxColLengths[i]-colLength)
		}

		fmt.Fprintln(f.out, s(f.indent*3)+"| "+strings.Join(cols, " | ")+" |")
	}
}

// longest gives a list of longest columns of all rows in Table
func maxColLengths(t *messages.PickleTable, clrs ...colors.ColorFunc) []int {
	if t == nil {
		return []int{}
	}

	longest := make([]int, len(t.Rows[0].Cells))
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			for _, c := range clrs {
				ln := utf8.RuneCountInString(c(cell.Value))
				if longest[i] < ln {
					longest[i] = ln
				}
			}

			ln := utf8.RuneCountInString(cell.Value)
			if longest[i] < ln {
				longest[i] = ln
			}
		}
	}

	return longest
}

func longestExampleRow(t *messages.Examples, clrs ...colors.ColorFunc) []int {
	if t == nil {
		return []int{}
	}

	longest := make([]int, len(t.TableHeader.Cells))
	for i, cell := range t.TableHeader.Cells {
		for _, c := range clrs {
			ln := utf8.RuneCountInString(c(cell.Value))
			if longest[i] < ln {
				longest[i] = ln
			}
		}

		ln := utf8.RuneCountInString(cell.Value)
		if longest[i] < ln {
			longest[i] = ln
		}
	}

	for _, row := range t.TableBody {
		for i, cell := range row.Cells {
			for _, c := range clrs {
				ln := utf8.RuneCountInString(c(cell.Value))
				if longest[i] < ln {
					longest[i] = ln
				}
			}

			ln := utf8.RuneCountInString(cell.Value)
			if longest[i] < ln {
				longest[i] = ln
			}
		}
	}

	return longest
}

func (f *Pretty) longestStep(steps []*messages.Step, pickleLength int) int {
	max := pickleLength

	for _, step := range steps {
		length := f.lengthPickleStep(step.Keyword, step.Text)
		if length > max {
			max = length
		}
	}

	return max
}

// a line number representation in feature file
func line(path string, loc *messages.Location) string {
	// Path can contain a line number already.
	// This line number has to be trimmed to avoid duplication.
	path = strings.TrimSuffix(path, fmt.Sprintf(":%d", loc.Line))
	return " " + blackb(fmt.Sprintf("# %s:%d", path, loc.Line))
}

func (f *Pretty) lengthPickleStep(keyword, text string) int {
	return f.indent*2 + utf8.RuneCountInString(strings.TrimSpace(keyword)+" "+text)
}

func (f *Pretty) lengthPickle(keyword, name string) int {
	return f.indent + utf8.RuneCountInString(strings.TrimSpace(keyword)+": "+name)
}
