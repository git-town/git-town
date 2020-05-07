package godog

/*
   The specification for the formatting originated from https://www.relishapp.com/cucumber/cucumber/docs/formatters/json-output-formatter.
   I found that documentation was misleading or out dated.  To validate formatting I create a ruby cucumber test harness and ran the
   same feature files through godog and the ruby cucumber.

   The docstrings in the cucumber.feature represent the cucumber output for those same feature definitions.

   I did note that comments in ruby could be at just about any level in particular Feature, Scenario and Step.  In godog I
   could only find comments under the Feature data structure.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cucumber/messages-go/v10"
)

func init() {
	Format("cucumber", "Produces cucumber JSON format output.", cucumberFunc)
}

func cucumberFunc(suite string, out io.Writer) Formatter {
	return &cukefmt{basefmt: newBaseFmt(suite, out)}
}

// Replace spaces with - This function is used to create the "id" fields of the cucumber output.
func makeID(name string) string {
	return strings.Replace(strings.ToLower(name), " ", "-", -1)
}

// The sequence of type structs are used to marshall the json object.
type cukeComment struct {
	Value string `json:"value"`
	Line  int    `json:"line"`
}

type cukeDocstring struct {
	Value       string `json:"value"`
	ContentType string `json:"content_type"`
	Line        int    `json:"line"`
}

type cukeTag struct {
	Name string `json:"name"`
	Line int    `json:"line"`
}

type cukeResult struct {
	Status   string `json:"status"`
	Error    string `json:"error_message,omitempty"`
	Duration *int   `json:"duration,omitempty"`
}

type cukeMatch struct {
	Location string `json:"location"`
}

type cukeStep struct {
	Keyword   string              `json:"keyword"`
	Name      string              `json:"name"`
	Line      int                 `json:"line"`
	Docstring *cukeDocstring      `json:"doc_string,omitempty"`
	Match     cukeMatch           `json:"match"`
	Result    cukeResult          `json:"result"`
	DataTable []*cukeDataTableRow `json:"rows,omitempty"`
}

type cukeDataTableRow struct {
	Cells []string `json:"cells"`
}

type cukeElement struct {
	ID          string     `json:"id"`
	Keyword     string     `json:"keyword"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Line        int        `json:"line"`
	Type        string     `json:"type"`
	Tags        []cukeTag  `json:"tags,omitempty"`
	Steps       []cukeStep `json:"steps,omitempty"`
}

type cukeFeatureJSON struct {
	URI         string        `json:"uri"`
	ID          string        `json:"id"`
	Keyword     string        `json:"keyword"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Line        int           `json:"line"`
	Comments    []cukeComment `json:"comments,omitempty"`
	Tags        []cukeTag     `json:"tags,omitempty"`
	Elements    []cukeElement `json:"elements,omitempty"`
}

type cukefmt struct {
	*basefmt

	// currently running feature path, to be part of id.
	// this is sadly not passed by gherkin nodes.
	// it restricts this formatter to run only in synchronous single
	// threaded execution. Unless running a copy of formatter for each feature
	path       string
	status     stepResultStatus  // last step status, before skipped
	ID         string            // current test id.
	results    []cukeFeatureJSON // structure that represent cuke results
	curStep    *cukeStep         // track the current step
	curElement *cukeElement      // track the current element
	curFeature *cukeFeatureJSON  // track the current feature
	curOutline cukeElement       // Each example show up as an outline element but the outline is parsed only once
	// so I need to keep track of the current outline
	curRow         int       // current row of the example table as it is being processed.
	curExampleTags []cukeTag // temporary storage for tags associate with the current example table.
	startTime      time.Time // used to time duration of the step execution
	curExampleName string    // Due to the fact that examples are parsed once and then iterated over for each result then we need to keep track
	// of the example name inorder to build id fields.
}

func (f *cukefmt) Pickle(pickle *messages.Pickle) {
	f.basefmt.Pickle(pickle)

	scenario := f.findScenario(pickle.AstNodeIds[0])

	f.curFeature.Elements = append(f.curFeature.Elements, cukeElement{})
	f.curElement = &f.curFeature.Elements[len(f.curFeature.Elements)-1]

	f.curElement.Name = pickle.Name
	f.curElement.Line = int(scenario.Location.Line)
	f.curElement.Description = scenario.Description
	f.curElement.Keyword = scenario.Keyword
	f.curElement.ID = f.curFeature.ID + ";" + makeID(pickle.Name)
	f.curElement.Type = "scenario"

	f.curElement.Tags = make([]cukeTag, len(scenario.Tags)+len(f.curFeature.Tags))

	if len(f.curElement.Tags) > 0 {
		// apply feature level tags
		copy(f.curElement.Tags, f.curFeature.Tags)

		// apply scenario level tags.
		for idx, element := range scenario.Tags {
			f.curElement.Tags[idx+len(f.curFeature.Tags)].Line = int(element.Location.Line)
			f.curElement.Tags[idx+len(f.curFeature.Tags)].Name = element.Name
		}
	}

	if len(pickle.AstNodeIds) == 1 {
		return
	}

	example, _ := f.findExample(pickle.AstNodeIds[1])
	// apply example level tags.
	for _, tag := range example.Tags {
		tag := cukeTag{Line: int(tag.Location.Line), Name: tag.Name}
		f.curElement.Tags = append(f.curElement.Tags, tag)
	}

	examples := scenario.GetExamples()
	if len(examples) > 0 {
		rowID := pickle.AstNodeIds[1]

		for _, example := range examples {
			for idx, row := range example.TableBody {
				if rowID == row.Id {
					f.curElement.ID += fmt.Sprintf(";%s;%d", makeID(example.Name), idx+2)
					f.curElement.Line = int(row.Location.Line)
				}
			}
		}
	}

}

func (f *cukefmt) Feature(gd *messages.GherkinDocument, p string, c []byte) {
	f.basefmt.Feature(gd, p, c)

	f.path = p
	f.ID = makeID(gd.Feature.Name)
	f.results = append(f.results, cukeFeatureJSON{})

	f.curFeature = &f.results[len(f.results)-1]
	f.curFeature.URI = p
	f.curFeature.Name = gd.Feature.Name
	f.curFeature.Keyword = gd.Feature.Keyword
	f.curFeature.Line = int(gd.Feature.Location.Line)
	f.curFeature.Description = gd.Feature.Description
	f.curFeature.ID = f.ID
	f.curFeature.Tags = make([]cukeTag, len(gd.Feature.Tags))

	for idx, element := range gd.Feature.Tags {
		f.curFeature.Tags[idx].Line = int(element.Location.Line)
		f.curFeature.Tags[idx].Name = element.Name
	}

	f.curFeature.Comments = make([]cukeComment, len(gd.Comments))
	for idx, comment := range gd.Comments {
		f.curFeature.Comments[idx].Value = strings.TrimSpace(comment.Text)
		f.curFeature.Comments[idx].Line = int(comment.Location.Line)
	}

}

func (f *cukefmt) Summary() {
	dat, err := json.MarshalIndent(f.results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(f.out, "%s\n", string(dat))
}

func (f *cukefmt) step(res *stepResult) {
	d := int(timeNowFunc().Sub(f.startTime).Nanoseconds())
	f.curStep.Result.Duration = &d
	f.curStep.Result.Status = res.status.String()
	if res.err != nil {
		f.curStep.Result.Error = res.err.Error()
	}
}

func (f *cukefmt) Defined(pickle *messages.Pickle, pickleStep *messages.Pickle_PickleStep, def *StepDefinition) {
	f.startTime = timeNowFunc() // start timing the step
	f.curElement.Steps = append(f.curElement.Steps, cukeStep{})
	f.curStep = &f.curElement.Steps[len(f.curElement.Steps)-1]

	step := f.findStep(pickleStep.AstNodeIds[0])

	line := step.Location.Line
	if len(pickle.AstNodeIds) == 2 {
		_, row := f.findExample(pickle.AstNodeIds[1])
		line = row.Location.Line
	}

	f.curStep.Name = pickleStep.Text
	f.curStep.Line = int(line)
	f.curStep.Keyword = step.Keyword

	arg := pickleStep.Argument

	if arg.GetDocString() != nil && step.GetDocString() != nil {
		f.curStep.Docstring = &cukeDocstring{}
		f.curStep.Docstring.ContentType = strings.TrimSpace(arg.GetDocString().MediaType)
		f.curStep.Docstring.Line = int(step.GetDocString().Location.Line)
		f.curStep.Docstring.Value = arg.GetDocString().Content
	}

	if arg.GetDataTable() != nil {
		f.curStep.DataTable = make([]*cukeDataTableRow, len(arg.GetDataTable().Rows))
		for i, row := range arg.GetDataTable().Rows {
			cells := make([]string, len(row.Cells))
			for j, cell := range row.Cells {
				cells[j] = cell.Value
			}
			f.curStep.DataTable[i] = &cukeDataTableRow{Cells: cells}
		}
	}

	if def != nil {
		f.curStep.Match.Location = strings.Split(def.definitionID(), " ")[0]
	}
}

func (f *cukefmt) Passed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Passed(pickle, step, match)

	f.status = passed
	f.step(f.lastStepResult())
}

func (f *cukefmt) Skipped(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Skipped(pickle, step, match)

	f.step(f.lastStepResult())

	// no duration reported for skipped.
	f.curStep.Result.Duration = nil
}

func (f *cukefmt) Undefined(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Undefined(pickle, step, match)

	f.status = undefined
	f.step(f.lastStepResult())

	// the location for undefined is the feature file location not the step file.
	f.curStep.Match.Location = fmt.Sprintf("%s:%d", f.path, f.findStep(step.AstNodeIds[0]).Location.Line)
	f.curStep.Result.Duration = nil
}

func (f *cukefmt) Failed(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition, err error) {
	f.basefmt.Failed(pickle, step, match, err)

	f.status = failed
	f.step(f.lastStepResult())
}

func (f *cukefmt) Pending(pickle *messages.Pickle, step *messages.Pickle_PickleStep, match *StepDefinition) {
	f.basefmt.Pending(pickle, step, match)

	f.status = pending
	f.step(f.lastStepResult())

	// the location for pending is the feature file location not the step file.
	f.curStep.Match.Location = fmt.Sprintf("%s:%d", f.path, f.findStep(step.AstNodeIds[0]).Location.Line)
	f.curStep.Result.Duration = nil
}
