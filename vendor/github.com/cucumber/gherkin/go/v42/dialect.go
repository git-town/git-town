package gherkin

import (
	"sort"

	messages "github.com/cucumber/messages/go/v34"
)

func init() {
	// only needs to be calcuated one, on startup
	for _, d := range builtinDialects {
		d.stepKeywords = []string{}
		d.stepKeywords = append(d.stepKeywords, d.Keywords[given]...)
		d.stepKeywords = append(d.stepKeywords, d.Keywords[when]...)
		d.stepKeywords = append(d.stepKeywords, d.Keywords[then]...)
		d.stepKeywords = append(d.stepKeywords, d.Keywords[and]...)
		d.stepKeywords = append(d.stepKeywords, d.Keywords[but]...)
		sort.Slice(d.stepKeywords, func(i, j int) bool {
			return d.stepKeywords[i] > d.stepKeywords[j]
		})
	}
}

type Dialect struct {
	Language     string
	Name         string
	Native       string
	Keywords     map[string][]string
	KeywordTypes map[string]messages.StepKeywordType

	stepKeywords []string
}

func (g *Dialect) FeatureKeywords() []string {
	return g.Keywords["feature"]
}

func (g *Dialect) RuleKeywords() []string {
	return g.Keywords["rule"]
}

func (g *Dialect) ScenarioKeywords() []string {
	return g.Keywords["scenario"]
}

func (g *Dialect) StepKeywords() []string {
	return g.stepKeywords
}

func (g *Dialect) BackgroundKeywords() []string {
	return g.Keywords["background"]
}

func (g *Dialect) ScenarioOutlineKeywords() []string {
	return g.Keywords["scenarioOutline"]
}

func (g *Dialect) ExamplesKeywords() []string {
	return g.Keywords["examples"]
}

func (g *Dialect) StepKeywordType(keyword string) messages.StepKeywordType {
	return g.KeywordTypes[keyword]
}

type DialectProvider interface {
	GetDialect(language string) *Dialect
}

type gherkinDialectMap map[string]*Dialect

func (g gherkinDialectMap) GetDialect(language string) *Dialect {
	return g[language]
}
