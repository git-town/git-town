package gherkin

import (
	"github.com/cucumber/messages-go/v16"
	"strings"
)

func Pickles(gherkinDocument messages.GherkinDocument, uri string, newId func() string) []*messages.Pickle {
	pickles := make([]*messages.Pickle, 0)
	if gherkinDocument.Feature == nil {
		return pickles
	}
	language := gherkinDocument.Feature.Language

	pickles = compileFeature(pickles, *gherkinDocument.Feature, uri, language, newId)
	return pickles
}

func compileFeature(pickles []*messages.Pickle, feature messages.Feature, uri string, language string, newId func() string) []*messages.Pickle {
	featureBackgroundSteps := make([]*messages.Step, 0)
	featureTags := feature.Tags
	for _, child := range feature.Children {
		if child.Background != nil {
			featureBackgroundSteps = append(featureBackgroundSteps, child.Background.Steps...)
		}
		if child.Rule != nil {
			pickles = compileRule(pickles, child.Rule, featureTags, featureBackgroundSteps, uri, language, newId)
		}
		if child.Scenario != nil {
			if len(child.Scenario.Examples) == 0 {
				pickles = compileScenario(pickles, featureBackgroundSteps, child.Scenario, featureTags, uri, language, newId)
			} else {
				pickles = compileScenarioOutline(pickles, child.Scenario, featureTags, featureBackgroundSteps, uri, language, newId)
			}
		}
	}
	return pickles
}

func compileRule(
	pickles []*messages.Pickle,
	rule *messages.Rule,
	featureTags []*messages.Tag,
	featureBackgroundSteps []*messages.Step,
	uri string,
	language string,
	newId func() string,
) []*messages.Pickle {
	ruleBackgroundSteps := make([]*messages.Step, 0)
	ruleBackgroundSteps = append(ruleBackgroundSteps, featureBackgroundSteps...)
	tags := append(featureTags, rule.Tags...)

	for _, child := range rule.Children {
		if child.Background != nil {
			ruleBackgroundSteps = append(ruleBackgroundSteps, child.Background.Steps...)
		}
		if child.Scenario != nil {
			if len(child.Scenario.Examples) == 0 {
				pickles = compileScenario(pickles, ruleBackgroundSteps, child.Scenario, tags, uri, language, newId)
			} else {
				pickles = compileScenarioOutline(pickles, child.Scenario, tags, ruleBackgroundSteps, uri, language, newId)
			}
		}
	}
	return pickles

}

func compileScenarioOutline(
	pickles []*messages.Pickle,
	scenario *messages.Scenario,
	inheritedTags []*messages.Tag,
	backgroundSteps []*messages.Step,
	uri string,
	language string,
	newId func() string,
) []*messages.Pickle {
	for _, examples := range scenario.Examples {
		if examples.TableHeader == nil {
			continue
		}
		variableCells := examples.TableHeader.Cells
		for _, valuesRow := range examples.TableBody {
			valueCells := valuesRow.Cells

			computedPickleSteps := make([]*messages.PickleStep, 0)
			pickleBackgroundSteps := make([]*messages.PickleStep, 0)

			if len(scenario.Steps) > 0 {
				pickleBackgroundSteps = pickleSteps(backgroundSteps, newId)
			}

			// translate computedPickleSteps based on valuesRow
			for _, step := range scenario.Steps {
				text := step.Text
				for i, variableCell := range variableCells {
					text = strings.Replace(text, "<"+variableCell.Value+">", valueCells[i].Value, -1)
				}

				pickleStep := pickleStep(step, variableCells, valuesRow, newId)
				computedPickleSteps = append(computedPickleSteps, pickleStep)
			}

			// translate pickle name
			name := scenario.Name
			for i, key := range variableCells {
				name = strings.Replace(name, "<"+key.Value+">", valueCells[i].Value, -1)
			}

			if len(computedPickleSteps) > 0 {
				computedPickleSteps = append(pickleBackgroundSteps, computedPickleSteps...)
			}

			id := newId()
			tags := pickleTags(append(inheritedTags, append(scenario.Tags, examples.Tags...)...))

			pickles = append(pickles, &messages.Pickle{
				Id:         id,
				Uri:        uri,
				Steps:      computedPickleSteps,
				Tags:       tags,
				Name:       name,
				Language:   language,
				AstNodeIds: []string{scenario.Id, valuesRow.Id},
			})
		}
	}
	return pickles
}

func compileScenario(
	pickles []*messages.Pickle,
	backgroundSteps []*messages.Step,
	scenario *messages.Scenario,
	inheritedTags []*messages.Tag,
	uri string,
	language string,
	newId func() string,
) []*messages.Pickle {
	steps := make([]*messages.PickleStep, 0)
	if len(scenario.Steps) > 0 {
		pickleBackgroundSteps := pickleSteps(backgroundSteps, newId)
		steps = append(pickleBackgroundSteps, pickleSteps(scenario.Steps, newId)...)
	}
	tags := pickleTags(append(inheritedTags, scenario.Tags...))
	id := newId()
	pickles = append(pickles, &messages.Pickle{
		Id:         id,
		Uri:        uri,
		Steps:      steps,
		Tags:       tags,
		Name:       scenario.Name,
		Language:   language,
		AstNodeIds: []string{scenario.Id},
	})
	return pickles
}

func pickleDataTable(table *messages.DataTable, variableCells []*messages.TableCell, valueCells []*messages.TableCell) *messages.PickleTable {
	pickleTableRows := make([]*messages.PickleTableRow, len(table.Rows))
	for i, row := range table.Rows {
		pickleTableCells := make([]*messages.PickleTableCell, len(row.Cells))
		for j, cell := range row.Cells {
			pickleTableCells[j] = &messages.PickleTableCell{
				Value: interpolate(cell.Value, variableCells, valueCells),
			}
		}
		pickleTableRows[i] = &messages.PickleTableRow{Cells: pickleTableCells}
	}
	return &messages.PickleTable{Rows: pickleTableRows}
}

func pickleDocString(docString *messages.DocString, variableCells []*messages.TableCell, valueCells []*messages.TableCell) *messages.PickleDocString {
	return &messages.PickleDocString{
		MediaType: interpolate(docString.MediaType, variableCells, valueCells),
		Content:   interpolate(docString.Content, variableCells, valueCells),
	}
}

func pickleTags(tags []*messages.Tag) []*messages.PickleTag {
	ptags := make([]*messages.PickleTag, len(tags))
	for i, tag := range tags {
		ptags[i] = &messages.PickleTag{
			Name:      tag.Name,
			AstNodeId: tag.Id,
		}
	}
	return ptags
}

func pickleSteps(steps []*messages.Step, newId func() string) []*messages.PickleStep {
	pickleSteps := make([]*messages.PickleStep, len(steps))
	for i, step := range steps {
		pickleStep := pickleStep(step, nil, nil, newId)
		pickleSteps[i] = pickleStep
	}
	return pickleSteps
}

func pickleStep(
	step *messages.Step,
	variableCells []*messages.TableCell,
	valuesRow *messages.TableRow,
	newId func() string) *messages.PickleStep {

	var valueCells []*messages.TableCell
	if valuesRow != nil {
		valueCells = valuesRow.Cells
	}

	pickleStep := &messages.PickleStep{
		Id:         newId(),
		Text:       interpolate(step.Text, variableCells, valueCells),
		AstNodeIds: []string{step.Id},
	}
	if valuesRow != nil {
		pickleStep.AstNodeIds = append(pickleStep.AstNodeIds, valuesRow.Id)
	}
	if step.DataTable != nil {
		pickleStep.Argument = &messages.PickleStepArgument{
			DataTable: pickleDataTable(step.DataTable, variableCells, valueCells),
		}
	}
	if step.DocString != nil {
		pickleStep.Argument = &messages.PickleStepArgument{
			DocString: pickleDocString(step.DocString, variableCells, valueCells),
		}
	}
	return pickleStep
}

func interpolate(s string, variableCells []*messages.TableCell, valueCells []*messages.TableCell) string {
	if variableCells == nil || valueCells == nil {
		return s
	}

	for i, variableCell := range variableCells {
		s = strings.Replace(s, "<"+variableCell.Value+">", valueCells[i].Value, -1)
	}

	return s
}
