package gherkin

import (
	"github.com/cucumber/messages-go/v10"
	"strings"
)

type AstBuilder interface {
	Builder
	GetGherkinDocument() *messages.GherkinDocument
}

type astBuilder struct {
	stack    []*astNode
	comments []*messages.GherkinDocument_Comment
	newId    func() string
}

func (t *astBuilder) Reset() {
	t.comments = []*messages.GherkinDocument_Comment{}
	t.stack = []*astNode{}
	t.push(newAstNode(RuleTypeNone))
}

func (t *astBuilder) GetGherkinDocument() *messages.GherkinDocument {
	res := t.currentNode().getSingle(RuleTypeGherkinDocument)
	if val, ok := res.(*messages.GherkinDocument); ok {
		return val
	}
	return nil
}

type astNode struct {
	ruleType RuleType
	subNodes map[RuleType][]interface{}
}

func (a *astNode) add(rt RuleType, obj interface{}) {
	a.subNodes[rt] = append(a.subNodes[rt], obj)
}

func (a *astNode) getSingle(rt RuleType) interface{} {
	if val, ok := a.subNodes[rt]; ok {
		for i := range val {
			return val[i]
		}
	}
	return nil
}

func (a *astNode) getItems(rt RuleType) []interface{} {
	var res []interface{}
	if val, ok := a.subNodes[rt]; ok {
		for i := range val {
			res = append(res, val[i])
		}
	}
	return res
}

func (a *astNode) getToken(tt TokenType) *Token {
	if val, ok := a.getSingle(tt.RuleType()).(*Token); ok {
		return val
	}
	return nil
}

func (a *astNode) getTokens(tt TokenType) []*Token {
	var items = a.getItems(tt.RuleType())
	var tokens []*Token
	for i := range items {
		if val, ok := items[i].(*Token); ok {
			tokens = append(tokens, val)
		}
	}
	return tokens
}

func (t *astBuilder) currentNode() *astNode {
	if len(t.stack) > 0 {
		return t.stack[len(t.stack)-1]
	}
	return nil
}

func newAstNode(rt RuleType) *astNode {
	return &astNode{
		ruleType: rt,
		subNodes: make(map[RuleType][]interface{}),
	}
}

func NewAstBuilder(newId func() string) AstBuilder {
	builder := new(astBuilder)
	builder.newId = newId
	builder.comments = []*messages.GherkinDocument_Comment{}
	builder.push(newAstNode(RuleTypeNone))
	return builder
}

func (t *astBuilder) push(n *astNode) {
	t.stack = append(t.stack, n)
}

func (t *astBuilder) pop() *astNode {
	x := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return x
}

func (t *astBuilder) Build(tok *Token) (bool, error) {
	if tok.Type == TokenTypeComment {
		comment := &messages.GherkinDocument_Comment{
			Location: astLocation(tok),
			Text:     tok.Text,
		}
		t.comments = append(t.comments, comment)
	} else {
		t.currentNode().add(tok.Type.RuleType(), tok)
	}
	return true, nil
}

func (t *astBuilder) StartRule(r RuleType) (bool, error) {
	t.push(newAstNode(r))
	return true, nil
}

func (t *astBuilder) EndRule(r RuleType) (bool, error) {
	node := t.pop()
	transformedNode, err := t.transformNode(node)
	t.currentNode().add(node.ruleType, transformedNode)
	return true, err
}

func (t *astBuilder) transformNode(node *astNode) (interface{}, error) {
	switch node.ruleType {

	case RuleTypeStep:
		stepLine := node.getToken(TokenTypeStepLine)

		step := &messages.GherkinDocument_Feature_Step{
			Location: astLocation(stepLine),
			Keyword:  stepLine.Keyword,
			Text:     stepLine.Text,
			Id:       t.newId(),
		}
		dataTable := node.getSingle(RuleTypeDataTable)
		if dataTable != nil {
			step.Argument = &messages.GherkinDocument_Feature_Step_DataTable_{
				DataTable: dataTable.(*messages.GherkinDocument_Feature_Step_DataTable),
			}
		} else {
			docString := node.getSingle(RuleTypeDocString)
			if docString != nil {
				step.Argument = &messages.GherkinDocument_Feature_Step_DocString_{DocString: docString.(*messages.GherkinDocument_Feature_Step_DocString)}
			}
		}

		return step, nil

	case RuleTypeDocString:
		separatorToken := node.getToken(TokenTypeDocStringSeparator)
		lineTokens := node.getTokens(TokenTypeOther)
		var text string
		for i := range lineTokens {
			if i > 0 {
				text += "\n"
			}
			text += lineTokens[i].Text
		}
		ds := &messages.GherkinDocument_Feature_Step_DocString{
			Location:  astLocation(separatorToken),
			MediaType: separatorToken.Text,
			Content:   text,
			Delimiter: separatorToken.Keyword,
		}
		return ds, nil

	case RuleTypeDataTable:
		rows, err := astTableRows(node, t.newId)
		dt := &messages.GherkinDocument_Feature_Step_DataTable{
			Location: rows[0].Location,
			Rows:     rows,
		}
		return dt, err

	case RuleTypeBackground:
		backgroundLine := node.getToken(TokenTypeBackgroundLine)
		description, _ := node.getSingle(RuleTypeDescription).(string)
		bg := &messages.GherkinDocument_Feature_Background{
			Location:    astLocation(backgroundLine),
			Keyword:     backgroundLine.Keyword,
			Name:        backgroundLine.Text,
			Description: description,
			Steps:       astSteps(node),
		}
		return bg, nil

	case RuleTypeScenarioDefinition:
		tags := astTags(node, t.newId)
		scenarioNode, _ := node.getSingle(RuleTypeScenario).(*astNode)

		scenarioLine := scenarioNode.getToken(TokenTypeScenarioLine)
		description, _ := scenarioNode.getSingle(RuleTypeDescription).(string)
		sc := &messages.GherkinDocument_Feature_Scenario{
			Id:          t.newId(),
			Tags:        tags,
			Location:    astLocation(scenarioLine),
			Keyword:     scenarioLine.Keyword,
			Name:        scenarioLine.Text,
			Description: description,
			Steps:       astSteps(scenarioNode),
			Examples:    astExamples(scenarioNode),
		}

		return sc, nil

	case RuleTypeExamplesDefinition:
		tags := astTags(node, t.newId)
		examplesNode, _ := node.getSingle(RuleTypeExamples).(*astNode)
		examplesLine := examplesNode.getToken(TokenTypeExamplesLine)
		description, _ := examplesNode.getSingle(RuleTypeDescription).(string)
		examplesTable := examplesNode.getSingle(RuleTypeExamplesTable)

		// TODO: Is this mutation style ok?
		ex := &messages.GherkinDocument_Feature_Scenario_Examples{}
		ex.Tags = tags
		ex.Location = astLocation(examplesLine)
		ex.Keyword = examplesLine.Keyword
		ex.Name = examplesLine.Text
		ex.Description = description
		ex.TableHeader = nil
		ex.TableBody = nil
		if examplesTable != nil {
			allRows, _ := examplesTable.([]*messages.GherkinDocument_Feature_TableRow)
			ex.TableHeader = allRows[0]
			ex.TableBody = allRows[1:]
		}
		return ex, nil

	case RuleTypeExamplesTable:
		allRows, err := astTableRows(node, t.newId)
		return allRows, err

	case RuleTypeDescription:
		lineTokens := node.getTokens(TokenTypeOther)
		// Trim trailing empty lines
		end := len(lineTokens)
		for end > 0 && strings.TrimSpace(lineTokens[end-1].Text) == "" {
			end--
		}
		var desc []string
		for i := range lineTokens[0:end] {
			desc = append(desc, lineTokens[i].Text)
		}
		return strings.Join(desc, "\n"), nil

	case RuleTypeFeature:
		header, ok := node.getSingle(RuleTypeFeatureHeader).(*astNode)
		if !ok {
			return nil, nil
		}
		tags := astTags(header, t.newId)
		featureLine := header.getToken(TokenTypeFeatureLine)
		if featureLine == nil {
			return nil, nil
		}

		var children []*messages.GherkinDocument_Feature_FeatureChild
		background, _ := node.getSingle(RuleTypeBackground).(*messages.GherkinDocument_Feature_Background)
		if background != nil {
			children = append(children, &messages.GherkinDocument_Feature_FeatureChild{
				Value: &messages.GherkinDocument_Feature_FeatureChild_Background{Background: background},
			})
		}
		scenarios := node.getItems(RuleTypeScenarioDefinition)
		for i := range scenarios {
			scenario := scenarios[i].(*messages.GherkinDocument_Feature_Scenario)
			children = append(children, &messages.GherkinDocument_Feature_FeatureChild{
				Value: &messages.GherkinDocument_Feature_FeatureChild_Scenario{Scenario: scenario},
			})
		}
		rules := node.getItems(RuleTypeRule)
		for i := range rules {
			rule := rules[i].(*messages.GherkinDocument_Feature_FeatureChild_Rule)
			children = append(children, &messages.GherkinDocument_Feature_FeatureChild{
				Value: &messages.GherkinDocument_Feature_FeatureChild_Rule_{
					Rule: rule,
				},
			})
		}

		description, _ := header.getSingle(RuleTypeDescription).(string)

		feat := &messages.GherkinDocument_Feature{}
		feat.Tags = tags
		feat.Location = astLocation(featureLine)
		feat.Language = featureLine.GherkinDialect
		feat.Keyword = featureLine.Keyword
		feat.Name = featureLine.Text
		feat.Description = description
		feat.Children = children
		return feat, nil

	case RuleTypeRule:
		header, ok := node.getSingle(RuleTypeRuleHeader).(*astNode)
		if !ok {
			return nil, nil
		}
		ruleLine := header.getToken(TokenTypeRuleLine)
		if ruleLine == nil {
			return nil, nil
		}

		var children []*messages.GherkinDocument_Feature_FeatureChild_RuleChild
		background, _ := node.getSingle(RuleTypeBackground).(*messages.GherkinDocument_Feature_Background)

		if background != nil {
			children = append(children, &messages.GherkinDocument_Feature_FeatureChild_RuleChild{
				Value: &messages.GherkinDocument_Feature_FeatureChild_RuleChild_Background{Background: background},
			})
		}
		scenarios := node.getItems(RuleTypeScenarioDefinition)
		for i := range scenarios {
			scenario := scenarios[i].(*messages.GherkinDocument_Feature_Scenario)
			children = append(children, &messages.GherkinDocument_Feature_FeatureChild_RuleChild{
				Value: &messages.GherkinDocument_Feature_FeatureChild_RuleChild_Scenario{Scenario: scenario},
			})
		}

		description, _ := header.getSingle(RuleTypeDescription).(string)

		rule := &messages.GherkinDocument_Feature_FeatureChild_Rule{}
		rule.Location = astLocation(ruleLine)
		rule.Keyword = ruleLine.Keyword
		rule.Name = ruleLine.Text
		rule.Description = description
		rule.Children = children
		return rule, nil

	case RuleTypeGherkinDocument:
		feature, _ := node.getSingle(RuleTypeFeature).(*messages.GherkinDocument_Feature)

		doc := &messages.GherkinDocument{}
		if feature != nil {
			doc.Feature = feature
		}
		doc.Comments = t.comments
		return doc, nil
	}
	return node, nil
}

func astLocation(t *Token) *messages.Location {
	return &messages.Location{
		Line:   uint32(t.Location.Line),
		Column: uint32(t.Location.Column),
	}
}

func astTableRows(t *astNode, newId func() string) (rows []*messages.GherkinDocument_Feature_TableRow, err error) {
	rows = []*messages.GherkinDocument_Feature_TableRow{}
	tokens := t.getTokens(TokenTypeTableRow)
	for i := range tokens {
		row := &messages.GherkinDocument_Feature_TableRow{
			Id:       newId(),
			Location: astLocation(tokens[i]),
			Cells:    astTableCells(tokens[i]),
		}
		rows = append(rows, row)
	}
	err = ensureCellCount(rows)
	return
}

func ensureCellCount(rows []*messages.GherkinDocument_Feature_TableRow) error {
	if len(rows) <= 1 {
		return nil
	}
	cellCount := len(rows[0].Cells)
	for i := range rows {
		if cellCount != len(rows[i].Cells) {
			return &parseError{"inconsistent cell count within the table", &Location{
				Line:   int(rows[i].Location.Line),
				Column: int(rows[i].Location.Column),
			}}
		}
	}
	return nil
}

func astTableCells(t *Token) (cells []*messages.GherkinDocument_Feature_TableRow_TableCell) {
	cells = []*messages.GherkinDocument_Feature_TableRow_TableCell{}
	for i := range t.Items {
		item := t.Items[i]
		cell := &messages.GherkinDocument_Feature_TableRow_TableCell{}
		cell.Location = &messages.Location{
			Line:   uint32(t.Location.Line),
			Column: uint32(item.Column),
		}
		cell.Value = item.Text
		cells = append(cells, cell)
	}
	return
}

func astSteps(t *astNode) (steps []*messages.GherkinDocument_Feature_Step) {
	steps = []*messages.GherkinDocument_Feature_Step{}
	tokens := t.getItems(RuleTypeStep)
	for i := range tokens {
		step, _ := tokens[i].(*messages.GherkinDocument_Feature_Step)
		steps = append(steps, step)
	}
	return
}

func astExamples(t *astNode) (examples []*messages.GherkinDocument_Feature_Scenario_Examples) {
	examples = []*messages.GherkinDocument_Feature_Scenario_Examples{}
	tokens := t.getItems(RuleTypeExamplesDefinition)
	for i := range tokens {
		example, _ := tokens[i].(*messages.GherkinDocument_Feature_Scenario_Examples)
		examples = append(examples, example)
	}
	return
}

func astTags(node *astNode, newId func() string) (tags []*messages.GherkinDocument_Feature_Tag) {
	tags = []*messages.GherkinDocument_Feature_Tag{}
	tagsNode, ok := node.getSingle(RuleTypeTags).(*astNode)
	if !ok {
		return
	}
	tokens := tagsNode.getTokens(TokenTypeTagLine)
	for i := range tokens {
		token := tokens[i]
		for k := range token.Items {
			item := token.Items[k]
			tag := &messages.GherkinDocument_Feature_Tag{}
			tag.Location = &messages.Location{
				Line:   uint32(token.Location.Line),
				Column: uint32(item.Column),
			}
			tag.Name = item.Text
			tag.Id = newId()
			tags = append(tags, tag)
		}
	}
	return
}
