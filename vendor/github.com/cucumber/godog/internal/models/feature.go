package models

import (
	messages "github.com/cucumber/messages/go/v21"
)

// Feature is an internal object to group together
// the parsed gherkin document, the pickles and the
// raw content.
type Feature struct {
	*messages.GherkinDocument
	Pickles []*messages.Pickle
	Content []byte
}

// FindRule returns the rule to which the given scenario belongs
func (f Feature) FindRule(astScenarioID string) *messages.Rule {
	for _, child := range f.GherkinDocument.Feature.Children {
		if ru := child.Rule; ru != nil {
			if rc := child.Rule; rc != nil {
				for _, rcc := range rc.Children {
					if sc := rcc.Scenario; sc != nil && sc.Id == astScenarioID {
						return ru
					}
				}
			}
		}
	}
	return nil
}

// FindScenario returns the scenario in the feature or in a rule in the feature
func (f Feature) FindScenario(astScenarioID string) *messages.Scenario {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.Scenario; sc != nil && sc.Id == astScenarioID {
			return sc
		}
		if rc := child.Rule; rc != nil {
			for _, rcc := range rc.Children {
				if sc := rcc.Scenario; sc != nil && sc.Id == astScenarioID {
					return sc
				}
			}
		}
	}

	return nil
}

// FindBackground ...
func (f Feature) FindBackground(astScenarioID string) *messages.Background {
	var bg *messages.Background

	for _, child := range f.GherkinDocument.Feature.Children {
		if tmp := child.Background; tmp != nil {
			bg = tmp
		}

		if sc := child.Scenario; sc != nil && sc.Id == astScenarioID {
			return bg
		}

		if ru := child.Rule; ru != nil {
			for _, rc := range ru.Children {
				if tmp := rc.Background; tmp != nil {
					bg = tmp
				}

				if sc := rc.Scenario; sc != nil && sc.Id == astScenarioID {
					return bg
				}
			}
		}
	}

	return nil
}

// FindExample ...
func (f Feature) FindExample(exampleAstID string) (*messages.Examples, *messages.TableRow) {
	for _, child := range f.GherkinDocument.Feature.Children {
		if sc := child.Scenario; sc != nil {
			for _, example := range sc.Examples {
				for _, row := range example.TableBody {
					if row.Id == exampleAstID {
						return example, row
					}
				}
			}
		}
		if ru := child.Rule; ru != nil {
			for _, rc := range ru.Children {
				if sc := rc.Scenario; sc != nil {
					for _, example := range sc.Examples {
						for _, row := range example.TableBody {
							if row.Id == exampleAstID {
								return example, row
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}

// FindStep ...
func (f Feature) FindStep(astStepID string) *messages.Step {
	for _, child := range f.GherkinDocument.Feature.Children {

		if ru := child.Rule; ru != nil {
			for _, ch := range ru.Children {
				if sc := ch.Scenario; sc != nil {
					for _, step := range sc.Steps {
						if step.Id == astStepID {
							return step
						}
					}
				}

				if bg := ch.Background; bg != nil {
					for _, step := range bg.Steps {
						if step.Id == astStepID {
							return step
						}
					}
				}
			}
		}

		if sc := child.Scenario; sc != nil {
			for _, step := range sc.Steps {
				if step.Id == astStepID {
					return step
				}
			}
		}

		if bg := child.Background; bg != nil {
			for _, step := range bg.Steps {
				if step.Id == astStepID {
					return step
				}
			}
		}
	}

	return nil
}
