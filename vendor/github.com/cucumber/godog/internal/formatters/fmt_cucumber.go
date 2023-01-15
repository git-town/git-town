package formatters

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
	"sort"
	"strings"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/formatters"
	"github.com/cucumber/godog/internal/models"
)

func init() {
	formatters.Format("cucumber", "Produces cucumber JSON format output.", CucumberFormatterFunc)
}

// CucumberFormatterFunc implements the FormatterFunc for the cucumber formatter
func CucumberFormatterFunc(suite string, out io.Writer) formatters.Formatter {
	return &Cuke{Base: NewBase(suite, out)}
}

// Cuke ...
type Cuke struct {
	*Base
}

// Summary renders test result as Cucumber JSON.
func (f *Cuke) Summary() {
	features := f.Storage.MustGetFeatures()

	res := f.buildCukeFeatures(features)

	dat, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(f.out, "%s\n", string(dat))
}

func (f *Cuke) buildCukeFeatures(features []*models.Feature) (res []CukeFeatureJSON) {
	sort.Sort(sortFeaturesByName(features))

	res = make([]CukeFeatureJSON, len(features))

	for idx, feat := range features {
		cukeFeature := buildCukeFeature(feat)

		pickles := f.Storage.MustGetPickles(feat.Uri)
		sort.Sort(sortPicklesByID(pickles))

		cukeFeature.Elements = f.buildCukeElements(pickles)

		for jdx, elem := range cukeFeature.Elements {
			elem.ID = cukeFeature.ID + ";" + makeCukeID(elem.Name) + elem.ID
			elem.Tags = append(cukeFeature.Tags, elem.Tags...)
			cukeFeature.Elements[jdx] = elem
		}

		res[idx] = cukeFeature
	}

	return res
}

func (f *Cuke) buildCukeElements(pickles []*messages.Pickle) (res []cukeElement) {
	res = make([]cukeElement, len(pickles))

	for idx, pickle := range pickles {
		pickleResult := f.Storage.MustGetPickleResult(pickle.Id)
		pickleStepResults := f.Storage.MustGetPickleStepResultsByPickleID(pickle.Id)

		cukeElement := f.buildCukeElement(pickle)

		stepStartedAt := pickleResult.StartedAt

		cukeElement.Steps = make([]cukeStep, len(pickleStepResults))
		sort.Sort(sortPickleStepResultsByPickleStepID(pickleStepResults))

		for jdx, stepResult := range pickleStepResults {
			cukeStep := f.buildCukeStep(pickle, stepResult)

			stepResultFinishedAt := stepResult.FinishedAt
			d := int(stepResultFinishedAt.Sub(stepStartedAt).Nanoseconds())
			stepStartedAt = stepResultFinishedAt

			cukeStep.Result.Duration = &d
			if stepResult.Status == undefined ||
				stepResult.Status == pending ||
				stepResult.Status == skipped {
				cukeStep.Result.Duration = nil
			}

			cukeElement.Steps[jdx] = cukeStep
		}

		res[idx] = cukeElement
	}

	return res
}

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

// CukeFeatureJSON ...
type CukeFeatureJSON struct {
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

func buildCukeFeature(feat *models.Feature) CukeFeatureJSON {
	cukeFeature := CukeFeatureJSON{
		URI:         feat.Uri,
		ID:          makeCukeID(feat.Feature.Name),
		Keyword:     feat.Feature.Keyword,
		Name:        feat.Feature.Name,
		Description: feat.Feature.Description,
		Line:        int(feat.Feature.Location.Line),
		Comments:    make([]cukeComment, len(feat.Comments)),
		Tags:        make([]cukeTag, len(feat.Feature.Tags)),
	}

	for idx, element := range feat.Feature.Tags {
		cukeFeature.Tags[idx].Line = int(element.Location.Line)
		cukeFeature.Tags[idx].Name = element.Name
	}

	for idx, comment := range feat.Comments {
		cukeFeature.Comments[idx].Value = strings.TrimSpace(comment.Text)
		cukeFeature.Comments[idx].Line = int(comment.Location.Line)
	}

	return cukeFeature
}

func (f *Cuke) buildCukeElement(pickle *messages.Pickle) (cukeElement cukeElement) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	scenario := feature.FindScenario(pickle.AstNodeIds[0])

	cukeElement.Name = pickle.Name
	cukeElement.Line = int(scenario.Location.Line)
	cukeElement.Description = scenario.Description
	cukeElement.Keyword = scenario.Keyword
	cukeElement.Type = "scenario"

	cukeElement.Tags = make([]cukeTag, len(scenario.Tags))
	for idx, element := range scenario.Tags {
		cukeElement.Tags[idx].Line = int(element.Location.Line)
		cukeElement.Tags[idx].Name = element.Name
	}

	if len(pickle.AstNodeIds) == 1 {
		return
	}

	example, _ := feature.FindExample(pickle.AstNodeIds[1])

	for _, tag := range example.Tags {
		tag := cukeTag{Line: int(tag.Location.Line), Name: tag.Name}
		cukeElement.Tags = append(cukeElement.Tags, tag)
	}

	examples := scenario.Examples
	if len(examples) > 0 {
		rowID := pickle.AstNodeIds[1]

		for _, example := range examples {
			for idx, row := range example.TableBody {
				if rowID == row.Id {
					cukeElement.ID += fmt.Sprintf(";%s;%d", makeCukeID(example.Name), idx+2)
					cukeElement.Line = int(row.Location.Line)
				}
			}
		}
	}

	return cukeElement
}

func (f *Cuke) buildCukeStep(pickle *messages.Pickle, stepResult models.PickleStepResult) (cukeStep cukeStep) {
	feature := f.Storage.MustGetFeature(pickle.Uri)
	pickleStep := f.Storage.MustGetPickleStep(stepResult.PickleStepID)
	step := feature.FindStep(pickleStep.AstNodeIds[0])

	line := step.Location.Line
	if len(pickle.AstNodeIds) == 2 {
		_, row := feature.FindExample(pickle.AstNodeIds[1])
		line = row.Location.Line
	}

	cukeStep.Name = pickleStep.Text
	cukeStep.Line = int(line)
	cukeStep.Keyword = step.Keyword

	arg := pickleStep.Argument

	if arg != nil {
		if arg.DocString != nil && step.DocString != nil {
			cukeStep.Docstring = &cukeDocstring{}
			cukeStep.Docstring.ContentType = strings.TrimSpace(arg.DocString.MediaType)
			if step.Location != nil {
				cukeStep.Docstring.Line = int(step.DocString.Location.Line)
			}
			cukeStep.Docstring.Value = arg.DocString.Content
		}

		if arg.DataTable != nil {
			cukeStep.DataTable = make([]*cukeDataTableRow, len(arg.DataTable.Rows))
			for i, row := range arg.DataTable.Rows {
				cells := make([]string, len(row.Cells))
				for j, cell := range row.Cells {
					cells[j] = cell.Value
				}
				cukeStep.DataTable[i] = &cukeDataTableRow{Cells: cells}
			}
		}
	}

	if stepResult.Def != nil {
		cukeStep.Match.Location = strings.Split(DefinitionID(stepResult.Def), " ")[0]
	}

	cukeStep.Result.Status = stepResult.Status.String()
	if stepResult.Err != nil {
		cukeStep.Result.Error = stepResult.Err.Error()
	}

	if stepResult.Status == undefined || stepResult.Status == pending {
		cukeStep.Match.Location = fmt.Sprintf("%s:%d", pickle.Uri, step.Location.Line)
	}

	return cukeStep
}

func makeCukeID(name string) string {
	return strings.Replace(strings.ToLower(name), " ", "-", -1)
}
