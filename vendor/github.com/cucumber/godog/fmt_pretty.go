package godog

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog/colors"
)

func init() {
	Format("pretty", "Prints every feature with runtime statuses.", prettyFunc)
}

func prettyFunc(suite string, out io.Writer) Formatter {
	return &pretty{basefmt: newBaseFmt(suite, out)}
}

var outlinePlaceholderRegexp = regexp.MustCompile("<[^>]+>")

// a built in default pretty formatter
type pretty struct {
	*basefmt
}

func (f *pretty) Feature(gd *messages.GherkinDocument, p string, c []byte) {
	f.basefmt.Feature(gd, p, c)
	f.printFeature(gd.Feature)
}

// Pickle takes a gherkin node for formatting
func (f *pretty) Pickle(pickle *messages.Pickle) {
	f.basefmt.Pickle(pickle)

	if len(pickle.Steps) == 0 {
		f.printUndefinedPickle(pickle)
		return
	}
}

func (f *pretty) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Passed(pickle, step, match)
	f.printStep(f.lastStepResult())
}

func (f *pretty) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Skipped(pickle, step, match)
	f.printStep(f.lastStepResult())
}

func (f *pretty) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Undefined(pickle, step, match)
	f.printStep(f.lastStepResult())
}

func (f *pretty) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
	f.basefmt.Failed(pickle, step, match, err)
	f.printStep(f.lastStepResult())
}

func (f *pretty) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Pending(pickle, step, match)
	f.printStep(f.lastStepResult())
}

func (f *pretty) printFeature(feature *messages.GherkinDocument_Feature) {
	if len(f.features) > 1 {
		fmt.Fprintln(f.out, "") // not a first feature, add a newline
	}

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

func (f *pretty) scenarioLengths(scenarioAstID string) (scenarioHeaderLength int, maxLength int) {
	astScenario := f.findScenario(scenarioAstID)
	astBackground := f.findBackground(scenarioAstID)

	scenarioHeaderLength = f.lengthPickle(astScenario.Keyword, astScenario.Name)
	maxLength = f.longestStep(astScenario.Steps, scenarioHeaderLength)

	if astBackground != nil {
		maxLength = f.longestStep(astBackground.Steps, maxLength)
	}

	return scenarioHeaderLength, maxLength
}

func (f *pretty) printScenarioHeader(astScenario *messages.GherkinDocument_Feature_Scenario, spaceFilling int) {
	text := s(f.indent) + keywordAndName(astScenario.Keyword, astScenario.Name)
	text += s(spaceFilling) + f.line(f.lastFeature().Path, astScenario.Location)
	fmt.Fprintln(f.out, "\n"+text)
}

func (f *pretty) printUndefinedPickle(pickle *messages.Pickle) {
	astScenario := f.findScenario(pickle.AstNodeIds[0])
	astBackground := f.findBackground(pickle.AstNodeIds[0])

	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle.AstNodeIds[0])

	if astBackground != nil {
		fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
		for _, step := range astBackground.Steps {
			text := s(f.indent*2) + cyan(strings.TrimSpace(step.Keyword)) + " " + cyan(step.Text)
			fmt.Fprintln(f.out, text)
		}
	}

	//  do not print scenario headers and examples multiple times
	if len(astScenario.Examples) > 0 {
		exampleTable, exampleRow := f.findExample(pickle.AstNodeIds[1])
		firstExampleRow := exampleTable.TableBody[0].Id == exampleRow.Id
		firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

		if !(firstExamplesTable && firstExampleRow) {
			return
		}
	}

	f.printScenarioHeader(astScenario, maxLength-scenarioHeaderLength)

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

// Summary sumarize the feature formatter output
func (f *pretty) Summary() {
	failedStepResults := f.findStepResults(failed)
	if len(failedStepResults) > 0 {
		fmt.Fprintln(f.out, "\n--- "+red("Failed steps:")+"\n")
		for _, fail := range failedStepResults {
			feature := f.findFeature(fail.owner.AstNodeIds[0])

			astScenario := f.findScenario(fail.owner.AstNodeIds[0])
			scenarioDesc := fmt.Sprintf("%s: %s", astScenario.Keyword, fail.owner.Name)

			astStep := f.findStep(fail.step.AstNodeIds[0])
			stepDesc := strings.TrimSpace(astStep.Keyword) + " " + fail.step.Text

			fmt.Fprintln(f.out, s(f.indent)+red(scenarioDesc)+f.line(feature.Path, astScenario.Location))
			fmt.Fprintln(f.out, s(f.indent*2)+red(stepDesc)+f.line(feature.Path, astStep.Location))
			fmt.Fprintln(f.out, s(f.indent*3)+red("Error: ")+redb(fmt.Sprintf("%+v", fail.err))+"\n")
		}
	}

	f.basefmt.Summary()
}

func (f *pretty) printOutlineExample(pickle *messages.Pickle, backgroundSteps int) {
	var errorMsg string
	var clr = green

	astScenario := f.findScenario(pickle.AstNodeIds[0])
	scenarioHeaderLength, maxLength := f.scenarioLengths(pickle.AstNodeIds[0])

	exampleTable, exampleRow := f.findExample(pickle.AstNodeIds[1])
	printExampleHeader := exampleTable.TableBody[0].Id == exampleRow.Id
	firstExamplesTable := astScenario.Examples[0].Location.Line == exampleTable.Location.Line

	firstExecutedScenarioStep := len(f.lastFeature().lastPickleResult().stepResults) == backgroundSteps+1
	if firstExamplesTable && printExampleHeader && firstExecutedScenarioStep {
		f.printScenarioHeader(astScenario, maxLength-scenarioHeaderLength)
	}

	if len(exampleTable.TableBody) == 0 {
		// do not print empty examples
		return
	}

	lastStep := len(f.lastFeature().lastPickleResult().stepResults) == len(pickle.Steps)
	if !lastStep {
		// do not print examples unless all steps has finished
		return
	}

	for _, result := range f.lastFeature().lastPickleResult().stepResults {
		// determine example row status
		switch {
		case result.status == failed:
			errorMsg = result.err.Error()
			clr = result.status.clr()
		case result.status == undefined || result.status == pending:
			clr = result.status.clr()
		case result.status == skipped && clr == nil:
			clr = cyan
		}

		if firstExamplesTable && printExampleHeader {
			// in first example, we need to print steps
			var text string

			astStep := f.findStep(result.step.AstNodeIds[0])

			if result.def != nil {
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

				_, maxLength := f.scenarioLengths(result.owner.AstNodeIds[0])
				stepLength := f.lengthPickleStep(astStep.Keyword, astStep.Text)

				text += s(maxLength - stepLength)
				text += " " + blackb("# "+result.def.definitionID())
			} else {
				text = cyan(astStep.Text)
			}
			// print the step outline
			fmt.Fprintln(f.out, s(f.indent*2)+cyan(strings.TrimSpace(astStep.Keyword))+" "+text)

			if table := result.step.Argument.GetDataTable(); table != nil {
				f.printTable(table, cyan)
			}

			if docString := astStep.GetDocString(); docString != nil {
				f.printDocString(docString)
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

func (f *pretty) printTableRow(row *messages.GherkinDocument_Feature_TableRow, max []int, clr colors.ColorFunc) {
	cells := make([]string, len(row.Cells))

	for i, cell := range row.Cells {
		val := clr(cell.Value)
		ln := utf8.RuneCountInString(val)
		cells[i] = val + s(max[i]-ln)
	}

	fmt.Fprintln(f.out, s(f.indent*3)+"| "+strings.Join(cells, " | ")+" |")
}

func (f *pretty) printTableHeader(row *messages.GherkinDocument_Feature_TableRow, max []int) {
	f.printTableRow(row, max, cyan)
}

func (f *pretty) printStep(result *stepResult) {
	astBackground := f.findBackground(result.owner.AstNodeIds[0])
	astScenario := f.findScenario(result.owner.AstNodeIds[0])
	astStep := f.findStep(result.step.AstNodeIds[0])

	var backgroundSteps int
	if astBackground != nil {
		backgroundSteps = len(astBackground.Steps)
	}

	astBackgroundStep := backgroundSteps > 0 && backgroundSteps >= len(f.lastFeature().lastPickleResult().stepResults)

	if astBackgroundStep {
		if len(f.lastFeature().pickleResults) > 1 {
			return
		}

		firstExecutedBackgroundStep := astBackground != nil && len(f.lastFeature().lastPickleResult().stepResults) == 1
		if firstExecutedBackgroundStep {
			fmt.Fprintln(f.out, "\n"+s(f.indent)+keywordAndName(astBackground.Keyword, astBackground.Name))
		}
	}

	if !astBackgroundStep && len(astScenario.Examples) > 0 {
		f.printOutlineExample(result.owner, backgroundSteps)
		return
	}

	scenarioHeaderLength, maxLength := f.scenarioLengths(result.owner.AstNodeIds[0])
	stepLength := f.lengthPickleStep(astStep.Keyword, result.step.Text)

	firstExecutedScenarioStep := len(f.lastFeature().lastPickleResult().stepResults) == backgroundSteps+1
	if !astBackgroundStep && firstExecutedScenarioStep {
		f.printScenarioHeader(astScenario, maxLength-scenarioHeaderLength)
	}

	text := s(f.indent*2) + result.status.clr()(strings.TrimSpace(astStep.Keyword)) + " " + result.status.clr()(result.step.Text)
	if result.def != nil {
		text += s(maxLength - stepLength + 1)
		text += blackb("# " + result.def.definitionID())
	}
	fmt.Fprintln(f.out, text)

	if table := result.step.Argument.GetDataTable(); table != nil {
		f.printTable(table, cyan)
	}

	if docString := astStep.GetDocString(); docString != nil {
		f.printDocString(docString)
	}

	if result.err != nil {
		fmt.Fprintln(f.out, s(f.indent*2)+redb(fmt.Sprintf("%+v", result.err)))
	}

	if result.status == pending {
		fmt.Fprintln(f.out, s(f.indent*3)+yellow("TODO: write pending definition"))
	}
}

func (f *pretty) printDocString(docString *messages.GherkinDocument_Feature_Step_DocString) {
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
func (f *pretty) printTable(t *messages.PickleStepArgument_PickleTable, c colors.ColorFunc) {
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
func maxColLengths(t *messages.PickleStepArgument_PickleTable, clrs ...colors.ColorFunc) []int {
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

func longestExampleRow(t *messages.GherkinDocument_Feature_Scenario_Examples, clrs ...colors.ColorFunc) []int {
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

func (f *pretty) longestStep(steps []*messages.GherkinDocument_Feature_Step, pickleLength int) int {
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
func (f *pretty) line(path string, loc *messages.Location) string {
	return " " + blackb(fmt.Sprintf("# %s:%d", path, loc.Line))
}

func (f *pretty) lengthPickleStep(keyword, text string) int {
	return f.indent*2 + utf8.RuneCountInString(strings.TrimSpace(keyword)+" "+text)
}

func (f *pretty) lengthPickle(keyword, name string) int {
	return f.indent + utf8.RuneCountInString(strings.TrimSpace(keyword)+": "+name)
}
