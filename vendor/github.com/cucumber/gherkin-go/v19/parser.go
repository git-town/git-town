//
// This file is generated. Do not edit! Edit parser.go.razor instead.

//
package gherkin

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenTypeNone TokenType = iota
	TokenTypeEOF
	TokenTypeEmpty
	TokenTypeComment
	TokenTypeTagLine
	TokenTypeFeatureLine
	TokenTypeRuleLine
	TokenTypeBackgroundLine
	TokenTypeScenarioLine
	TokenTypeExamplesLine
	TokenTypeStepLine
	TokenTypeDocStringSeparator
	TokenTypeTableRow
	TokenTypeLanguage
	TokenTypeOther
)

func tokenTypeForRule(rt RuleType) TokenType {
	return TokenTypeNone
}

func (t TokenType) Name() string {
	switch t {
	case TokenTypeEOF:
		return "EOF"
	case TokenTypeEmpty:
		return "Empty"
	case TokenTypeComment:
		return "Comment"
	case TokenTypeTagLine:
		return "TagLine"
	case TokenTypeFeatureLine:
		return "FeatureLine"
	case TokenTypeRuleLine:
		return "RuleLine"
	case TokenTypeBackgroundLine:
		return "BackgroundLine"
	case TokenTypeScenarioLine:
		return "ScenarioLine"
	case TokenTypeExamplesLine:
		return "ExamplesLine"
	case TokenTypeStepLine:
		return "StepLine"
	case TokenTypeDocStringSeparator:
		return "DocStringSeparator"
	case TokenTypeTableRow:
		return "TableRow"
	case TokenTypeLanguage:
		return "Language"
	case TokenTypeOther:
		return "Other"
	}
	return ""
}

func (t TokenType) RuleType() RuleType {
	switch t {
	case TokenTypeEOF:
		return RuleTypeEOF
	case TokenTypeEmpty:
		return RuleTypeEmpty
	case TokenTypeComment:
		return RuleTypeComment
	case TokenTypeTagLine:
		return RuleTypeTagLine
	case TokenTypeFeatureLine:
		return RuleTypeFeatureLine
	case TokenTypeRuleLine:
		return RuleTypeRuleLine
	case TokenTypeBackgroundLine:
		return RuleTypeBackgroundLine
	case TokenTypeScenarioLine:
		return RuleTypeScenarioLine
	case TokenTypeExamplesLine:
		return RuleTypeExamplesLine
	case TokenTypeStepLine:
		return RuleTypeStepLine
	case TokenTypeDocStringSeparator:
		return RuleTypeDocStringSeparator
	case TokenTypeTableRow:
		return RuleTypeTableRow
	case TokenTypeLanguage:
		return RuleTypeLanguage
	case TokenTypeOther:
		return RuleTypeOther
	}
	return RuleTypeNone
}

type RuleType int

const (
	RuleTypeNone RuleType = iota

	RuleTypeEOF
	RuleTypeEmpty
	RuleTypeComment
	RuleTypeTagLine
	RuleTypeFeatureLine
	RuleTypeRuleLine
	RuleTypeBackgroundLine
	RuleTypeScenarioLine
	RuleTypeExamplesLine
	RuleTypeStepLine
	RuleTypeDocStringSeparator
	RuleTypeTableRow
	RuleTypeLanguage
	RuleTypeOther
	RuleTypeGherkinDocument
	RuleTypeFeature
	RuleTypeFeatureHeader
	RuleTypeRule
	RuleTypeRuleHeader
	RuleTypeBackground
	RuleTypeScenarioDefinition
	RuleTypeScenario
	RuleTypeExamplesDefinition
	RuleTypeExamples
	RuleTypeExamplesTable
	RuleTypeStep
	RuleTypeStepArg
	RuleTypeDataTable
	RuleTypeDocString
	RuleTypeTags
	RuleTypeDescriptionHelper
	RuleTypeDescription
)

func (t RuleType) IsEOF() bool {
	return t == RuleTypeEOF
}
func (t RuleType) Name() string {
	switch t {
	case RuleTypeEOF:
		return "#EOF"
	case RuleTypeEmpty:
		return "#Empty"
	case RuleTypeComment:
		return "#Comment"
	case RuleTypeTagLine:
		return "#TagLine"
	case RuleTypeFeatureLine:
		return "#FeatureLine"
	case RuleTypeRuleLine:
		return "#RuleLine"
	case RuleTypeBackgroundLine:
		return "#BackgroundLine"
	case RuleTypeScenarioLine:
		return "#ScenarioLine"
	case RuleTypeExamplesLine:
		return "#ExamplesLine"
	case RuleTypeStepLine:
		return "#StepLine"
	case RuleTypeDocStringSeparator:
		return "#DocStringSeparator"
	case RuleTypeTableRow:
		return "#TableRow"
	case RuleTypeLanguage:
		return "#Language"
	case RuleTypeOther:
		return "#Other"
	case RuleTypeGherkinDocument:
		return "GherkinDocument"
	case RuleTypeFeature:
		return "Feature"
	case RuleTypeFeatureHeader:
		return "FeatureHeader"
	case RuleTypeRule:
		return "Rule"
	case RuleTypeRuleHeader:
		return "RuleHeader"
	case RuleTypeBackground:
		return "Background"
	case RuleTypeScenarioDefinition:
		return "ScenarioDefinition"
	case RuleTypeScenario:
		return "Scenario"
	case RuleTypeExamplesDefinition:
		return "ExamplesDefinition"
	case RuleTypeExamples:
		return "Examples"
	case RuleTypeExamplesTable:
		return "ExamplesTable"
	case RuleTypeStep:
		return "Step"
	case RuleTypeStepArg:
		return "StepArg"
	case RuleTypeDataTable:
		return "DataTable"
	case RuleTypeDocString:
		return "DocString"
	case RuleTypeTags:
		return "Tags"
	case RuleTypeDescriptionHelper:
		return "DescriptionHelper"
	case RuleTypeDescription:
		return "Description"
	}
	return ""
}

type Location struct {
	Line   int
	Column int
}

type parseError struct {
	msg string
	loc *Location
}

func (a *parseError) Error() string {
	return fmt.Sprintf("(%d:%d): %s", a.loc.Line, a.loc.Column, a.msg)
}

type parseErrors []error

func (pe parseErrors) Error() string {
	var ret = []string{"Parser errors:"}
	for i := range pe {
		ret = append(ret, pe[i].Error())
	}
	return strings.Join(ret, "\n")
}

func (p *parser) Parse(s Scanner, m Matcher) (err error) {
	p.builder.Reset()
	m.Reset()
	ctxt := &parseContext{p, s, p.builder, m, nil, nil}
	var state int
	ctxt.startRule(RuleTypeGherkinDocument)
	for {
		gl, eof, err := ctxt.scan()
		if err != nil {
			ctxt.addError(err)
			if p.stopAtFirstError {
				break
			}
		}
		state, err = ctxt.match(state, gl)
		if err != nil {
			ctxt.addError(err)
			if p.stopAtFirstError {
				break
			}
		}
		if eof {
			// done! \o/
			break
		}
	}
	ctxt.endRule(RuleTypeGherkinDocument)
	if len(ctxt.errors) > 0 {
		return ctxt.errors
	}
	return
}

type parseContext struct {
	p      *parser
	s      Scanner
	b      Builder
	m      Matcher
	queue  []*scanResult
	errors parseErrors
}

func (ctxt *parseContext) addError(e error) {
	ctxt.errors = append(ctxt.errors, e)
	// if (p.errors.length > 10)
	//   throw Errors.CompositeParserException.create(p.errors);
}

type scanResult struct {
	line  *Line
	atEof bool
	err   error
}

func (ctxt *parseContext) scan() (*Line, bool, error) {
	l := len(ctxt.queue)
	if l > 0 {
		x := ctxt.queue[0]
		ctxt.queue = ctxt.queue[1:]
		return x.line, x.atEof, x.err
	}
	return ctxt.s.Scan()
}

func (ctxt *parseContext) startRule(r RuleType) (bool, error) {
	ok, err := ctxt.b.StartRule(r)
	if err != nil {
		ctxt.addError(err)
	}
	return ok, err
}

func (ctxt *parseContext) endRule(r RuleType) (bool, error) {
	ok, err := ctxt.b.EndRule(r)
	if err != nil {
		ctxt.addError(err)
	}
	return ok, err
}

func (ctxt *parseContext) build(t *Token) (bool, error) {
	ok, err := ctxt.b.Build(t)
	if err != nil {
		ctxt.addError(err)
	}
	return ok, err
}

func (ctxt *parseContext) match(state int, line *Line) (newState int, err error) {
	switch state {
	case 0:
		return ctxt.matchAt0(line)
	case 1:
		return ctxt.matchAt1(line)
	case 2:
		return ctxt.matchAt2(line)
	case 3:
		return ctxt.matchAt3(line)
	case 4:
		return ctxt.matchAt4(line)
	case 5:
		return ctxt.matchAt5(line)
	case 6:
		return ctxt.matchAt6(line)
	case 7:
		return ctxt.matchAt7(line)
	case 8:
		return ctxt.matchAt8(line)
	case 9:
		return ctxt.matchAt9(line)
	case 10:
		return ctxt.matchAt10(line)
	case 11:
		return ctxt.matchAt11(line)
	case 12:
		return ctxt.matchAt12(line)
	case 13:
		return ctxt.matchAt13(line)
	case 14:
		return ctxt.matchAt14(line)
	case 15:
		return ctxt.matchAt15(line)
	case 16:
		return ctxt.matchAt16(line)
	case 17:
		return ctxt.matchAt17(line)
	case 18:
		return ctxt.matchAt18(line)
	case 19:
		return ctxt.matchAt19(line)
	case 20:
		return ctxt.matchAt20(line)
	case 21:
		return ctxt.matchAt21(line)
	case 22:
		return ctxt.matchAt22(line)
	case 23:
		return ctxt.matchAt23(line)
	case 24:
		return ctxt.matchAt24(line)
	case 25:
		return ctxt.matchAt25(line)
	case 26:
		return ctxt.matchAt26(line)
	case 27:
		return ctxt.matchAt27(line)
	case 28:
		return ctxt.matchAt28(line)
	case 29:
		return ctxt.matchAt29(line)
	case 30:
		return ctxt.matchAt30(line)
	case 31:
		return ctxt.matchAt31(line)
	case 32:
		return ctxt.matchAt32(line)
	case 33:
		return ctxt.matchAt33(line)
	case 34:
		return ctxt.matchAt34(line)
	case 35:
		return ctxt.matchAt35(line)
	case 36:
		return ctxt.matchAt36(line)
	case 37:
		return ctxt.matchAt37(line)
	case 38:
		return ctxt.matchAt38(line)
	case 39:
		return ctxt.matchAt39(line)
	case 40:
		return ctxt.matchAt40(line)
	case 41:
		return ctxt.matchAt41(line)
	case 43:
		return ctxt.matchAt43(line)
	case 44:
		return ctxt.matchAt44(line)
	case 45:
		return ctxt.matchAt45(line)
	case 46:
		return ctxt.matchAt46(line)
	case 47:
		return ctxt.matchAt47(line)
	case 48:
		return ctxt.matchAt48(line)
	case 49:
		return ctxt.matchAt49(line)
	case 50:
		return ctxt.matchAt50(line)
	default:
		return state, fmt.Errorf("Unknown state: %+v", state)
	}
}

// Start
func (ctxt *parseContext) matchAt0(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchLanguage(line); ok {
		ctxt.startRule(RuleTypeFeature)
		ctxt.startRule(RuleTypeFeatureHeader)
		ctxt.build(token)
		return 1, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.startRule(RuleTypeFeature)
		ctxt.startRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 2, err
	}
	if ok, token, err := ctxt.matchFeatureLine(line); ok {
		ctxt.startRule(RuleTypeFeature)
		ctxt.startRule(RuleTypeFeatureHeader)
		ctxt.build(token)
		return 3, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 0, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 0, err
	}

	// var stateComment = "State: 0 - Start"
	var expectedTokens = []string{"#EOF", "#Language", "#TagLine", "#FeatureLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 0, err
}

// GherkinDocument:0>Feature:0>FeatureHeader:0>#Language:0
func (ctxt *parseContext) matchAt1(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 2, err
	}
	if ok, token, err := ctxt.matchFeatureLine(line); ok {
		ctxt.build(token)
		return 3, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 1, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 1, err
	}

	// var stateComment = "State: 1 - GherkinDocument:0>Feature:0>FeatureHeader:0>#Language:0"
	var expectedTokens = []string{"#TagLine", "#FeatureLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 1, err
}

// GherkinDocument:0>Feature:0>FeatureHeader:1>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt2(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 2, err
	}
	if ok, token, err := ctxt.matchFeatureLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.build(token)
		return 3, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 2, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 2, err
	}

	// var stateComment = "State: 2 - GherkinDocument:0>Feature:0>FeatureHeader:1>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#FeatureLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 2, err
}

// GherkinDocument:0>Feature:0>FeatureHeader:2>#FeatureLine:0
func (ctxt *parseContext) matchAt3(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 3, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 5, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 6, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeFeatureHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 4, err
	}

	// var stateComment = "State: 3 - GherkinDocument:0>Feature:0>FeatureHeader:2>#FeatureLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 3, err
}

// GherkinDocument:0>Feature:0>FeatureHeader:3>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt4(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 5, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 6, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeFeatureHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 4, err
	}

	// var stateComment = "State: 4 - GherkinDocument:0>Feature:0>FeatureHeader:3>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 4, err
}

// GherkinDocument:0>Feature:0>FeatureHeader:3>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt5(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 5, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 6, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeFeatureHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeFeatureHeader)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 5, err
	}

	// var stateComment = "State: 5 - GherkinDocument:0>Feature:0>FeatureHeader:3>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 5, err
}

// GherkinDocument:0>Feature:1>Background:0>#BackgroundLine:0
func (ctxt *parseContext) matchAt6(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 6, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 8, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 7, err
	}

	// var stateComment = "State: 6 - GherkinDocument:0>Feature:1>Background:0>#BackgroundLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 6, err
}

// GherkinDocument:0>Feature:1>Background:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt7(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 8, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 7, err
	}

	// var stateComment = "State: 7 - GherkinDocument:0>Feature:1>Background:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 7, err
}

// GherkinDocument:0>Feature:1>Background:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt8(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 8, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 8, err
	}

	// var stateComment = "State: 8 - GherkinDocument:0>Feature:1>Background:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 8, err
}

// GherkinDocument:0>Feature:1>Background:2>Step:0>#StepLine:0
func (ctxt *parseContext) matchAt9(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeDataTable)
		ctxt.build(token)
		return 10, err
	}
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.startRule(RuleTypeDocString)
		ctxt.build(token)
		return 49, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 9, err
	}

	// var stateComment = "State: 9 - GherkinDocument:0>Feature:1>Background:2>Step:0>#StepLine:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#DocStringSeparator", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 9, err
}

// GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0
func (ctxt *parseContext) matchAt10(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 10, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 10, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 10, err
	}

	// var stateComment = "State: 10 - GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 10, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:0>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt11(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 11, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 11, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 11, err
	}

	// var stateComment = "State: 11 - GherkinDocument:0>Feature:2>ScenarioDefinition:0>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#ScenarioLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 11, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:0>#ScenarioLine:0
func (ctxt *parseContext) matchAt12(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 14, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 13, err
	}

	// var stateComment = "State: 12 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:0>#ScenarioLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 12, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt13(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 14, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 13, err
	}

	// var stateComment = "State: 13 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 13, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt14(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 14, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 14, err
	}

	// var stateComment = "State: 14 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 14, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:0>#StepLine:0
func (ctxt *parseContext) matchAt15(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeDataTable)
		ctxt.build(token)
		return 16, err
	}
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.startRule(RuleTypeDocString)
		ctxt.build(token)
		return 47, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 15, err
	}

	// var stateComment = "State: 15 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:0>#StepLine:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#DocStringSeparator", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 15, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0
func (ctxt *parseContext) matchAt16(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 16, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 16, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 16, err
	}

	// var stateComment = "State: 16 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 16, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:0>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt17(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 17, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 17, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 17, err
	}

	// var stateComment = "State: 17 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:0>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#ExamplesLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 17, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:0>#ExamplesLine:0
func (ctxt *parseContext) matchAt18(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 20, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 21, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 19, err
	}

	// var stateComment = "State: 18 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:0>#ExamplesLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 18, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt19(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 20, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 21, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 19, err
	}

	// var stateComment = "State: 19 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 19, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt20(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 20, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 21, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 20, err
	}

	// var stateComment = "State: 20 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 20, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:2>ExamplesTable:0>#TableRow:0
func (ctxt *parseContext) matchAt21(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 21, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamplesTable)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamplesTable)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 21, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 21, err
	}

	// var stateComment = "State: 21 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:2>ExamplesTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 21, err
}

// GherkinDocument:0>Feature:3>Rule:0>RuleHeader:0>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt22(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 22, err
	}

	// var stateComment = "State: 22 - GherkinDocument:0>Feature:3>Rule:0>RuleHeader:0>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 22, err
}

// GherkinDocument:0>Feature:3>Rule:0>RuleHeader:1>#RuleLine:0
func (ctxt *parseContext) matchAt23(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 25, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 26, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeRuleHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 24, err
	}

	// var stateComment = "State: 23 - GherkinDocument:0>Feature:3>Rule:0>RuleHeader:1>#RuleLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 23, err
}

// GherkinDocument:0>Feature:3>Rule:0>RuleHeader:2>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt24(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 25, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 26, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeRuleHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 24, err
	}

	// var stateComment = "State: 24 - GherkinDocument:0>Feature:3>Rule:0>RuleHeader:2>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 24, err
}

// GherkinDocument:0>Feature:3>Rule:0>RuleHeader:2>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt25(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 25, err
	}
	if ok, token, err := ctxt.matchBackgroundLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeBackground)
		ctxt.build(token)
		return 26, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeRuleHeader)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeRuleHeader)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 25, err
	}

	// var stateComment = "State: 25 - GherkinDocument:0>Feature:3>Rule:0>RuleHeader:2>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#BackgroundLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 25, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:0>#BackgroundLine:0
func (ctxt *parseContext) matchAt26(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 26, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 28, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 27, err
	}

	// var stateComment = "State: 26 - GherkinDocument:0>Feature:3>Rule:1>Background:0>#BackgroundLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 26, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt27(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 28, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 27, err
	}

	// var stateComment = "State: 27 - GherkinDocument:0>Feature:3>Rule:1>Background:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 27, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt28(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 28, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 28, err
	}

	// var stateComment = "State: 28 - GherkinDocument:0>Feature:3>Rule:1>Background:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 28, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:0>#StepLine:0
func (ctxt *parseContext) matchAt29(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeDataTable)
		ctxt.build(token)
		return 30, err
	}
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.startRule(RuleTypeDocString)
		ctxt.build(token)
		return 45, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 29, err
	}

	// var stateComment = "State: 29 - GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:0>#StepLine:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#DocStringSeparator", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 29, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0
func (ctxt *parseContext) matchAt30(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 30, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 30, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 30, err
	}

	// var stateComment = "State: 30 - GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 30, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:0>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt31(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 31, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 31, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 31, err
	}

	// var stateComment = "State: 31 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:0>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#ScenarioLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 31, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:0>#ScenarioLine:0
func (ctxt *parseContext) matchAt32(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 34, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 33, err
	}

	// var stateComment = "State: 32 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:0>#ScenarioLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 32, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt33(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 34, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 33, err
	}

	// var stateComment = "State: 33 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 33, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt34(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 34, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 34, err
	}

	// var stateComment = "State: 34 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 34, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:0>#StepLine:0
func (ctxt *parseContext) matchAt35(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeDataTable)
		ctxt.build(token)
		return 36, err
	}
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.startRule(RuleTypeDocString)
		ctxt.build(token)
		return 43, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 35, err
	}

	// var stateComment = "State: 35 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:0>#StepLine:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#DocStringSeparator", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 35, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0
func (ctxt *parseContext) matchAt36(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 36, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDataTable)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDataTable)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 36, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 36, err
	}

	// var stateComment = "State: 36 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:0>DataTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 36, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:0>Tags:0>#TagLine:0
func (ctxt *parseContext) matchAt37(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.build(token)
		return 37, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeTags)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 37, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 37, err
	}

	// var stateComment = "State: 37 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:0>Tags:0>#TagLine:0"
	var expectedTokens = []string{"#TagLine", "#ExamplesLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 37, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:0>#ExamplesLine:0
func (ctxt *parseContext) matchAt38(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 40, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 41, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.startRule(RuleTypeDescription)
		ctxt.build(token)
		return 39, err
	}

	// var stateComment = "State: 38 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:0>#ExamplesLine:0"
	var expectedTokens = []string{"#EOF", "#Empty", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 38, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:1>Description:0>#Other:0
func (ctxt *parseContext) matchAt39(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.build(token)
		return 40, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 41, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDescription)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDescription)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 39, err
	}

	// var stateComment = "State: 39 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:1>Description:0>#Other:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 39, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:2>#Comment:0
func (ctxt *parseContext) matchAt40(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 40, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.startRule(RuleTypeExamplesTable)
		ctxt.build(token)
		return 41, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 40, err
	}

	// var stateComment = "State: 40 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:1>DescriptionHelper:2>#Comment:0"
	var expectedTokens = []string{"#EOF", "#Comment", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 40, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:2>ExamplesTable:0>#TableRow:0
func (ctxt *parseContext) matchAt41(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchTableRow(line); ok {
		ctxt.build(token)
		return 41, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeExamplesTable)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeExamplesTable)
			ctxt.endRule(RuleTypeExamples)
			ctxt.endRule(RuleTypeExamplesDefinition)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeExamplesTable)
		ctxt.endRule(RuleTypeExamples)
		ctxt.endRule(RuleTypeExamplesDefinition)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 41, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 41, err
	}

	// var stateComment = "State: 41 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:3>ExamplesDefinition:1>Examples:2>ExamplesTable:0>#TableRow:0"
	var expectedTokens = []string{"#EOF", "#TableRow", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 41, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0
func (ctxt *parseContext) matchAt43(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.build(token)
		return 44, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 43, err
	}

	// var stateComment = "State: 43 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0"
	var expectedTokens = []string{"#DocStringSeparator", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 43, err
}

// GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0
func (ctxt *parseContext) matchAt44(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 35, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 37, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 38, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 44, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 44, err
	}

	// var stateComment = "State: 44 - GherkinDocument:0>Feature:3>Rule:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0"
	var expectedTokens = []string{"#EOF", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 44, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0
func (ctxt *parseContext) matchAt45(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.build(token)
		return 46, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 45, err
	}

	// var stateComment = "State: 45 - GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0"
	var expectedTokens = []string{"#DocStringSeparator", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 45, err
}

// GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0
func (ctxt *parseContext) matchAt46(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 29, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 31, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 32, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 46, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 46, err
	}

	// var stateComment = "State: 46 - GherkinDocument:0>Feature:3>Rule:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0"
	var expectedTokens = []string{"#EOF", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 46, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0
func (ctxt *parseContext) matchAt47(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.build(token)
		return 48, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 47, err
	}

	// var stateComment = "State: 47 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0"
	var expectedTokens = []string{"#DocStringSeparator", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 47, err
}

// GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0
func (ctxt *parseContext) matchAt48(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 15, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead1(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.startRule(RuleTypeExamplesDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 17, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeScenario)
			ctxt.endRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchExamplesLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeExamplesDefinition)
		ctxt.startRule(RuleTypeExamples)
		ctxt.build(token)
		return 18, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeScenario)
		ctxt.endRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 48, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 48, err
	}

	// var stateComment = "State: 48 - GherkinDocument:0>Feature:2>ScenarioDefinition:1>Scenario:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0"
	var expectedTokens = []string{"#EOF", "#StepLine", "#TagLine", "#ExamplesLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 48, err
}

// GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0
func (ctxt *parseContext) matchAt49(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchDocStringSeparator(line); ok {
		ctxt.build(token)
		return 50, err
	}
	if ok, token, err := ctxt.matchOther(line); ok {
		ctxt.build(token)
		return 49, err
	}

	// var stateComment = "State: 49 - GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:0>#DocStringSeparator:0"
	var expectedTokens = []string{"#DocStringSeparator", "#Other"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 49, err
}

// GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0
func (ctxt *parseContext) matchAt50(line *Line) (newState int, err error) {
	if ok, token, err := ctxt.matchEOF(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.endRule(RuleTypeFeature)
		ctxt.build(token)
		return 42, err
	}
	if ok, token, err := ctxt.matchStepLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.startRule(RuleTypeStep)
		ctxt.build(token)
		return 9, err
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		if ctxt.lookahead0(line) {
			ctxt.endRule(RuleTypeDocString)
			ctxt.endRule(RuleTypeStep)
			ctxt.endRule(RuleTypeBackground)
			ctxt.startRule(RuleTypeScenarioDefinition)
			ctxt.startRule(RuleTypeTags)
			ctxt.build(token)
			return 11, err
		}
	}
	if ok, token, err := ctxt.matchTagLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.startRule(RuleTypeTags)
		ctxt.build(token)
		return 22, err
	}
	if ok, token, err := ctxt.matchScenarioLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeScenarioDefinition)
		ctxt.startRule(RuleTypeScenario)
		ctxt.build(token)
		return 12, err
	}
	if ok, token, err := ctxt.matchRuleLine(line); ok {
		ctxt.endRule(RuleTypeDocString)
		ctxt.endRule(RuleTypeStep)
		ctxt.endRule(RuleTypeBackground)
		ctxt.startRule(RuleTypeRule)
		ctxt.startRule(RuleTypeRuleHeader)
		ctxt.build(token)
		return 23, err
	}
	if ok, token, err := ctxt.matchComment(line); ok {
		ctxt.build(token)
		return 50, err
	}
	if ok, token, err := ctxt.matchEmpty(line); ok {
		ctxt.build(token)
		return 50, err
	}

	// var stateComment = "State: 50 - GherkinDocument:0>Feature:1>Background:2>Step:1>StepArg:0>__alt0:1>DocString:2>#DocStringSeparator:0"
	var expectedTokens = []string{"#EOF", "#StepLine", "#TagLine", "#ScenarioLine", "#RuleLine", "#Comment", "#Empty"}
	if line.IsEof() {
		err = &parseError{
			msg: fmt.Sprintf("unexpected end of file, expected: %s", strings.Join(expectedTokens, ", ")),
			loc: &Location{Line: line.LineNumber, Column: 0},
		}
	} else {
		err = &parseError{
			msg: fmt.Sprintf("expected: %s, got '%s'", strings.Join(expectedTokens, ", "), line.LineText),
			loc: &Location{Line: line.LineNumber, Column: line.Indent() + 1},
		}
	}
	// if (ctxt.p.stopAtFirstError) throw error;
	//ctxt.addError(err)
	return 50, err
}

type Matcher interface {
	MatchEOF(line *Line) (bool, *Token, error)
	MatchEmpty(line *Line) (bool, *Token, error)
	MatchComment(line *Line) (bool, *Token, error)
	MatchTagLine(line *Line) (bool, *Token, error)
	MatchFeatureLine(line *Line) (bool, *Token, error)
	MatchRuleLine(line *Line) (bool, *Token, error)
	MatchBackgroundLine(line *Line) (bool, *Token, error)
	MatchScenarioLine(line *Line) (bool, *Token, error)
	MatchExamplesLine(line *Line) (bool, *Token, error)
	MatchStepLine(line *Line) (bool, *Token, error)
	MatchDocStringSeparator(line *Line) (bool, *Token, error)
	MatchTableRow(line *Line) (bool, *Token, error)
	MatchLanguage(line *Line) (bool, *Token, error)
	MatchOther(line *Line) (bool, *Token, error)
	Reset()
}

func (ctxt *parseContext) isMatchEOF(line *Line) bool {
	ok, _, _ := ctxt.matchEOF(line)
	return ok
}
func (ctxt *parseContext) matchEOF(line *Line) (bool, *Token, error) {
	return ctxt.m.MatchEOF(line)
}

func (ctxt *parseContext) isMatchEmpty(line *Line) bool {
	ok, _, _ := ctxt.matchEmpty(line)
	return ok
}
func (ctxt *parseContext) matchEmpty(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchEmpty(line)
}

func (ctxt *parseContext) isMatchComment(line *Line) bool {
	ok, _, _ := ctxt.matchComment(line)
	return ok
}
func (ctxt *parseContext) matchComment(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchComment(line)
}

func (ctxt *parseContext) isMatchTagLine(line *Line) bool {
	ok, _, _ := ctxt.matchTagLine(line)
	return ok
}
func (ctxt *parseContext) matchTagLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchTagLine(line)
}

func (ctxt *parseContext) isMatchFeatureLine(line *Line) bool {
	ok, _, _ := ctxt.matchFeatureLine(line)
	return ok
}
func (ctxt *parseContext) matchFeatureLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchFeatureLine(line)
}

func (ctxt *parseContext) isMatchRuleLine(line *Line) bool {
	ok, _, _ := ctxt.matchRuleLine(line)
	return ok
}
func (ctxt *parseContext) matchRuleLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchRuleLine(line)
}

func (ctxt *parseContext) isMatchBackgroundLine(line *Line) bool {
	ok, _, _ := ctxt.matchBackgroundLine(line)
	return ok
}
func (ctxt *parseContext) matchBackgroundLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchBackgroundLine(line)
}

func (ctxt *parseContext) isMatchScenarioLine(line *Line) bool {
	ok, _, _ := ctxt.matchScenarioLine(line)
	return ok
}
func (ctxt *parseContext) matchScenarioLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchScenarioLine(line)
}

func (ctxt *parseContext) isMatchExamplesLine(line *Line) bool {
	ok, _, _ := ctxt.matchExamplesLine(line)
	return ok
}
func (ctxt *parseContext) matchExamplesLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchExamplesLine(line)
}

func (ctxt *parseContext) isMatchStepLine(line *Line) bool {
	ok, _, _ := ctxt.matchStepLine(line)
	return ok
}
func (ctxt *parseContext) matchStepLine(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchStepLine(line)
}

func (ctxt *parseContext) isMatchDocStringSeparator(line *Line) bool {
	ok, _, _ := ctxt.matchDocStringSeparator(line)
	return ok
}
func (ctxt *parseContext) matchDocStringSeparator(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchDocStringSeparator(line)
}

func (ctxt *parseContext) isMatchTableRow(line *Line) bool {
	ok, _, _ := ctxt.matchTableRow(line)
	return ok
}
func (ctxt *parseContext) matchTableRow(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchTableRow(line)
}

func (ctxt *parseContext) isMatchLanguage(line *Line) bool {
	ok, _, _ := ctxt.matchLanguage(line)
	return ok
}
func (ctxt *parseContext) matchLanguage(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchLanguage(line)
}

func (ctxt *parseContext) isMatchOther(line *Line) bool {
	ok, _, _ := ctxt.matchOther(line)
	return ok
}
func (ctxt *parseContext) matchOther(line *Line) (bool, *Token, error) {
	if line.IsEof() {
		return false, nil, nil
	}
	return ctxt.m.MatchOther(line)
}

func (ctxt *parseContext) lookahead0(initialLine *Line) bool {
	var queue []*scanResult
	var match bool

	for {
		line, atEof, err := ctxt.scan()
		queue = append(queue, &scanResult{line, atEof, err})

		if false || ctxt.isMatchScenarioLine(line) {
			match = true
			break
		}
		if !(false || ctxt.isMatchEmpty(line) || ctxt.isMatchComment(line) || ctxt.isMatchTagLine(line)) {
			break
		}
		if atEof {
			break
		}
	}

	ctxt.queue = append(ctxt.queue, queue...)

	return match
}

func (ctxt *parseContext) lookahead1(initialLine *Line) bool {
	var queue []*scanResult
	var match bool

	for {
		line, atEof, err := ctxt.scan()
		queue = append(queue, &scanResult{line, atEof, err})

		if false || ctxt.isMatchExamplesLine(line) {
			match = true
			break
		}
		if !(false || ctxt.isMatchEmpty(line) || ctxt.isMatchComment(line) || ctxt.isMatchTagLine(line)) {
			break
		}
		if atEof {
			break
		}
	}

	ctxt.queue = append(ctxt.queue, queue...)

	return match
}
