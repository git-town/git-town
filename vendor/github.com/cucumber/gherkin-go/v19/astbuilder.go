package gherkin

import (
	"github.com/cucumber/messages-go/v16"
	"strings"
)

type AstBuilder interface {
	Builder
	GetGherkinDocument() *messages.GherkinDocument
}

type astBuilder struct {
	stack    []*astNode
	comments []*messages.Comment
	newId    func() string
}

func (t *astBuilder) Reset() {
	t.comments = []*messages.Comment{}
	t.stack = []*astNode{}
	t.push(newAstNode(RuleTypeNone))
}

func (t *astBuilder) GetGherkinDocument() *messages.GherkinDocument {
	res := t.currentNode().getSingle(RuleTypeGherkinDocument, nil)
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

func (a *astNode) getSingle(rt RuleType, defaultValue interface{}) interface{} {
	if val, ok := a.subNodes[rt]; ok {
		for i := range val {
			return val[i]
		}
	}
	return defaultValue
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
	if val, ok := a.getSingle(tt.RuleType(), nil).(*Token); ok {
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
	builder.comments = []*messages.Comment{}
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
		comment := &messages.Comment{
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

		step := &messages.Step{
			Location: astLocation(stepLine),
			Keyword:  stepLine.Keyword,
			Text:     stepLine.Text,
			Id:       t.newId(),
		}
		dataTable := node.getSingle(RuleTypeDataTable, nil)
		if dataTable != nil {
			step.DataTable = dataTable.(*messages.DataTable)
		} else {
			docString := node.getSingle(RuleTypeDocString, nil)
			if docString != nil {
				step.DocString = docString.(*messages.DocString)
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
		docString := &messages.DocString{
			Location:  astLocation(separatorToken),
			Content:   text,
			Delimiter: separatorToken.Keyword,
		}
		if len(separatorToken.Text) > 0 {
			docString.MediaType = separatorToken.Text
		}

		return docString, nil

	case RuleTypeDataTable:
		rows, err := astTableRows(node, t.newId)
		dt := &messages.DataTable{
			Location: rows[0].Location,
			Rows:     rows,
		}
		return dt, err

	case RuleTypeBackground:
		backgroundLine := node.getToken(TokenTypeBackgroundLine)
		bg := &messages.Background{
			Id:          t.newId(),
			Location:    astLocation(backgroundLine),
			Keyword:     backgroundLine.Keyword,
			Name:        backgroundLine.Text,
			Description: getDescription(node),
			Steps:       astSteps(node),
		}
		return bg, nil

	case RuleTypeScenarioDefinition:
		scenarioNode := node.getSingle(RuleTypeScenario, nil).(*astNode)
		scenarioLine := scenarioNode.getToken(TokenTypeScenarioLine)
		tags := astTags(node, t.newId)
		sc := &messages.Scenario{
			Id:          t.newId(),
			Tags:        tags,
			Location:    astLocation(scenarioLine),
			Keyword:     scenarioLine.Keyword,
			Name:        scenarioLine.Text,
			Description: getDescription(scenarioNode),
			Steps:       astSteps(scenarioNode),
			Examples:    astExamples(scenarioNode),
		}

		return sc, nil

	case RuleTypeExamplesDefinition:
		tags := astTags(node, t.newId)
		examplesNode := node.getSingle(RuleTypeExamples, nil).(*astNode)
		examplesLine := examplesNode.getToken(TokenTypeExamplesLine)
		examplesTable := examplesNode.getSingle(RuleTypeExamplesTable, make([]*messages.TableRow, 0)).([]*messages.TableRow)

		var tableHeader *messages.TableRow
		var tableBody []*messages.TableRow

		if len(examplesTable) > 0 {
			tableHeader = examplesTable[0]
			tableBody = examplesTable[1:]
		} else {
			tableHeader = nil
			tableBody = examplesTable
		}

		ex := &messages.Examples{
			Id:          t.newId(),
			Tags:        tags,
			Location:    astLocation(examplesLine),
			Keyword:     examplesLine.Keyword,
			Name:        examplesLine.Text,
			Description: getDescription(examplesNode),
			TableHeader: tableHeader,
			TableBody:   tableBody,
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
		header := node.getSingle(RuleTypeFeatureHeader, nil).(*astNode)
		tags := astTags(header, t.newId)
		featureLine := header.getToken(TokenTypeFeatureLine)
		if featureLine == nil {
			return nil, nil
		}

		children := make([]*messages.FeatureChild, 0)
		background, _ := node.getSingle(RuleTypeBackground, nil).(*messages.Background)
		if background != nil {
			children = append(children, &messages.FeatureChild{
				Background: background,
			})
		}
		scenarios := node.getItems(RuleTypeScenarioDefinition)
		for i := range scenarios {
			scenario := scenarios[i].(*messages.Scenario)
			children = append(children, &messages.FeatureChild{
				Scenario: scenario,
			})
		}
		rules := node.getItems(RuleTypeRule)
		for i := range rules {
			rule := rules[i].(*messages.Rule)
			children = append(children, &messages.FeatureChild{
				Rule: rule,
			})
		}

		feature := &messages.Feature{
			Tags:        tags,
			Location:    astLocation(featureLine),
			Language:    featureLine.GherkinDialect,
			Keyword:     featureLine.Keyword,
			Name:        featureLine.Text,
			Description: getDescription(header),
			Children:    children,
		}
		return feature, nil

	case RuleTypeRule:
		header := node.getSingle(RuleTypeRuleHeader, nil).(*astNode)
		ruleLine := header.getToken(TokenTypeRuleLine)
		if ruleLine == nil {
			return nil, nil
		}

		tags := astTags(header, t.newId)
		var children []*messages.RuleChild
		background, _ := node.getSingle(RuleTypeBackground, nil).(*messages.Background)

		if background != nil {
			children = append(children, &messages.RuleChild{
				Background: background,
			})
		}
		scenarios := node.getItems(RuleTypeScenarioDefinition)
		for i := range scenarios {
			scenario := scenarios[i].(*messages.Scenario)
			children = append(children, &messages.RuleChild{
				Scenario: scenario,
			})
		}

		rule := &messages.Rule{
			Id:          t.newId(),
			Location:    astLocation(ruleLine),
			Keyword:     ruleLine.Keyword,
			Name:        ruleLine.Text,
			Description: getDescription(header),
			Children:    children,
			Tags:        tags,
		}
		return rule, nil

	case RuleTypeGherkinDocument:
		feature, _ := node.getSingle(RuleTypeFeature, nil).(*messages.Feature)

		doc := &messages.GherkinDocument{}
		if feature != nil {
			doc.Feature = feature
		}
		doc.Comments = t.comments
		return doc, nil
	}
	return node, nil
}

func getDescription(node *astNode) string {
	return node.getSingle(RuleTypeDescription, "").(string)
}

func astLocation(t *Token) *messages.Location {
	return &messages.Location{
		Line:   int64(t.Location.Line),
		Column: int64(t.Location.Column),
	}
}

func astTableRows(t *astNode, newId func() string) (rows []*messages.TableRow, err error) {
	rows = []*messages.TableRow{}
	tokens := t.getTokens(TokenTypeTableRow)
	for i := range tokens {
		row := &messages.TableRow{
			Id:       newId(),
			Location: astLocation(tokens[i]),
			Cells:    astTableCells(tokens[i]),
		}
		rows = append(rows, row)
	}
	err = ensureCellCount(rows)
	return
}

func ensureCellCount(rows []*messages.TableRow) error {
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

func astTableCells(t *Token) (cells []*messages.TableCell) {
	cells = []*messages.TableCell{}
	for i := range t.Items {
		item := t.Items[i]
		cell := &messages.TableCell{}
		cell.Location = &messages.Location{
			Line:   int64(t.Location.Line),
			Column: int64(item.Column),
		}
		cell.Value = item.Text
		cells = append(cells, cell)
	}
	return
}

func astSteps(t *astNode) (steps []*messages.Step) {
	steps = []*messages.Step{}
	tokens := t.getItems(RuleTypeStep)
	for i := range tokens {
		step, _ := tokens[i].(*messages.Step)
		steps = append(steps, step)
	}
	return
}

func astExamples(t *astNode) (examples []*messages.Examples) {
	examples = []*messages.Examples{}
	tokens := t.getItems(RuleTypeExamplesDefinition)
	for i := range tokens {
		example, _ := tokens[i].(*messages.Examples)
		examples = append(examples, example)
	}
	return
}

func astTags(node *astNode, newId func() string) (tags []*messages.Tag) {
	tags = []*messages.Tag{}
	tagsNode, ok := node.getSingle(RuleTypeTags, nil).(*astNode)
	if !ok {
		return
	}
	tokens := tagsNode.getTokens(TokenTypeTagLine)
	for i := range tokens {
		token := tokens[i]
		for k := range token.Items {
			item := token.Items[k]
			tag := &messages.Tag{}
			tag.Location = &messages.Location{
				Line:   int64(token.Location.Line),
				Column: int64(item.Column),
			}
			tag.Name = item.Text
			tag.Id = newId()
			tags = append(tags, tag)
		}
	}
	return
}
