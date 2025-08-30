package gherkin

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	DefaultDialect                = "en"
	CommentPrefix                 = "#"
	TagPrefix                     = "@"
	TitleKeywordSeparator         = ":"
	TableCellSeparator            = '|'
	EscapeChar                    = '\\'
	EscapedNewline                = 'n'
	DocstringSeparator            = "\"\"\""
	DocstringAlternativeSeparator = "```"
)

type matcher struct {
	gdp                      DialectProvider
	defaultLang              string
	lang                     string
	dialect                  *Dialect
	activeDocStringSeparator string
	indentToRemove           int
	languagePattern          *regexp.Regexp
}

func NewMatcher(gdp DialectProvider) Matcher {
	return &matcher{
		gdp:             gdp,
		defaultLang:     DefaultDialect,
		lang:            DefaultDialect,
		dialect:         gdp.GetDialect(DefaultDialect),
		languagePattern: regexp.MustCompile("^\\s*#\\s*language\\s*:\\s*([a-zA-Z\\-_]+)\\s*$"),
	}
}

func NewLanguageMatcher(gdp DialectProvider, language string) Matcher {
	return &matcher{
		gdp:             gdp,
		defaultLang:     language,
		lang:            language,
		dialect:         gdp.GetDialect(language),
		languagePattern: regexp.MustCompile("^\\s*#\\s*language\\s*:\\s*([a-zA-Z\\-_]+)\\s*$"),
	}
}

func (m *matcher) Reset() {
	m.indentToRemove = 0
	m.activeDocStringSeparator = ""
	if m.lang != "en" {
		m.dialect = m.gdp.GetDialect(m.defaultLang)
		m.lang = "en"
	}
}

func (m *matcher) newTokenAtLocation(line, index int) (token *Token) {
	column := index + 1
	token = new(Token)
	token.GherkinDialect = m.lang
	token.Location = &Location{line, column}
	return
}

func (m *matcher) MatchEOF(line *Line) (ok bool, token *Token, err error) {
	if line.IsEof() {
		token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
		token.Type = TokenTypeEOF
	}
	return
}

func (m *matcher) MatchEmpty(line *Line) (ok bool, token *Token, err error) {
	if line.IsEmpty() {
		token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
		token.Type = TokenTypeEmpty
	}
	return
}

func (m *matcher) MatchComment(line *Line) (ok bool, token *Token, err error) {
	if line.StartsWith(CommentPrefix) {
		token, ok = m.newTokenAtLocation(line.LineNumber, 0), true
		token.Type = TokenTypeComment
		token.Text = line.LineText
	}
	return
}

func (m *matcher) MatchTagLine(line *Line) (ok bool, token *Token, err error) {
	if !line.StartsWith(TagPrefix) {
		return
	}
	commentDelimiter := regexp.MustCompile(`\s+` + CommentPrefix)
	uncommentedLine := commentDelimiter.Split(line.TrimmedLineText, 2)[0]
	var tags []*LineSpan
	var column = line.Indent() + 1

	splits := strings.Split(uncommentedLine, TagPrefix)
	for i := range splits {
		txt := strings.TrimRightFunc(splits[i], func(r rune) bool {
			return unicode.IsSpace(r)
		})
		if len(txt) == 0 {
			continue
		}
		if !regexp.MustCompile(`^\S+$`).MatchString(txt) {
			location := &Location{line.LineNumber, column}
			msg := "A tag may not contain whitespace"
			err = &parseError{msg, location}
			break
		}
		tags = append(tags, &LineSpan{column, TagPrefix + txt})
		column = column + utf8.RuneCountInString(splits[i]) + 1
	}

	token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
	token.Type = TokenTypeTagLine
	token.Items = tags

	return
}

func (m *matcher) matchTitleLine(line *Line, tokenType TokenType, keywords []string) (ok bool, token *Token, err error) {
	for i := range keywords {
		keyword := keywords[i]
		if line.StartsWith(keyword + TitleKeywordSeparator) {
			token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
			token.Type = tokenType
			token.Keyword = keyword
			token.Text = strings.Trim(line.TrimmedLineText[len(keyword)+1:], " ")
			return
		}
	}
	return
}

func (m *matcher) MatchFeatureLine(line *Line) (ok bool, token *Token, err error) {
	return m.matchTitleLine(line, TokenTypeFeatureLine, m.dialect.FeatureKeywords())
}
func (m *matcher) MatchRuleLine(line *Line) (ok bool, token *Token, err error) {
	return m.matchTitleLine(line, TokenTypeRuleLine, m.dialect.RuleKeywords())
}
func (m *matcher) MatchBackgroundLine(line *Line) (ok bool, token *Token, err error) {
	return m.matchTitleLine(line, TokenTypeBackgroundLine, m.dialect.BackgroundKeywords())
}
func (m *matcher) MatchScenarioLine(line *Line) (ok bool, token *Token, err error) {
	ok, token, err = m.matchTitleLine(line, TokenTypeScenarioLine, m.dialect.ScenarioKeywords())
	if ok || (err != nil) {
		return ok, token, err
	}
	ok, token, err = m.matchTitleLine(line, TokenTypeScenarioLine, m.dialect.ScenarioOutlineKeywords())
	return ok, token, err
}
func (m *matcher) MatchExamplesLine(line *Line) (ok bool, token *Token, err error) {
	return m.matchTitleLine(line, TokenTypeExamplesLine, m.dialect.ExamplesKeywords())
}
func (m *matcher) MatchStepLine(line *Line) (ok bool, token *Token, err error) {
	keywords := m.dialect.StepKeywords()
	for i := range keywords {
		keyword := keywords[i]
		if line.StartsWith(keyword) {
			token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
			token.Type = TokenTypeStepLine
			token.Keyword = keyword
			token.KeywordType = m.dialect.StepKeywordType(keyword)
			token.Text = strings.Trim(line.TrimmedLineText[len(keyword):], " ")
			return
		}
	}
	return
}

func (m *matcher) MatchDocStringSeparator(line *Line) (ok bool, token *Token, err error) {
	if m.activeDocStringSeparator != "" {
		if line.StartsWith(m.activeDocStringSeparator) {
			// close
			token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
			token.Type = TokenTypeDocStringSeparator
			token.Keyword = m.activeDocStringSeparator

			m.indentToRemove = 0
			m.activeDocStringSeparator = ""
		}
		return
	}
	if line.StartsWith(DocstringSeparator) {
		m.activeDocStringSeparator = DocstringSeparator
	} else if line.StartsWith(DocstringAlternativeSeparator) {
		m.activeDocStringSeparator = DocstringAlternativeSeparator
	}
	if m.activeDocStringSeparator != "" {
		// open
		mediaType := line.TrimmedLineText[len(m.activeDocStringSeparator):]
		m.indentToRemove = line.Indent()
		token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
		token.Type = TokenTypeDocStringSeparator
		token.Keyword = m.activeDocStringSeparator
		token.Text = mediaType
	}
	return
}

func isSpaceAndNotNewLine(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}

func (m *matcher) MatchTableRow(line *Line) (ok bool, token *Token, err error) {
	var firstChar, firstPos = utf8.DecodeRuneInString(line.TrimmedLineText)
	if firstChar == TableCellSeparator {
		var cells []*LineSpan
		var cell []rune
		var startCol = line.Indent() + 2 // column where the current cell started
		// start after the first separator, it's not included in the cell
		for i, w, col := firstPos, 0, startCol; i < len(line.TrimmedLineText); i += w {
			var char rune
			char, w = utf8.DecodeRuneInString(line.TrimmedLineText[i:])
			if char == TableCellSeparator {
				// append current cell
				txt := string(cell)

				txtTrimmedLeadingSpace := strings.TrimLeftFunc(txt, isSpaceAndNotNewLine)
				ind := utf8.RuneCountInString(txt) - utf8.RuneCountInString(txtTrimmedLeadingSpace)
				txtTrimmed := strings.TrimRightFunc(txtTrimmedLeadingSpace, isSpaceAndNotNewLine)
				cells = append(cells, &LineSpan{startCol + ind, txtTrimmed})
				// start building next
				cell = make([]rune, 0)
				startCol = col + 1
			} else if char == EscapeChar {
				// skip this character but count the column
				i += w
				col++
				char, w = utf8.DecodeRuneInString(line.TrimmedLineText[i:])
				if char == EscapedNewline {
					cell = append(cell, '\n')
				} else {
					if char != TableCellSeparator && char != EscapeChar {
						cell = append(cell, EscapeChar)
					}
					cell = append(cell, char)
				}
			} else {
				cell = append(cell, char)
			}
			col++
		}

		token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
		token.Type = TokenTypeTableRow
		token.Items = cells
	}
	return
}

func (m *matcher) MatchLanguage(line *Line) (ok bool, token *Token, err error) {
	matches := m.languagePattern.FindStringSubmatch(line.TrimmedLineText)
	if len(matches) > 0 {
		lang := matches[1]
		token, ok = m.newTokenAtLocation(line.LineNumber, line.Indent()), true
		token.Type = TokenTypeLanguage
		token.Text = lang

		dialect := m.gdp.GetDialect(lang)
		if dialect == nil {
			err = &parseError{"Language not supported: " + lang, token.Location}
		} else {
			m.lang = lang
			m.dialect = dialect
		}
	}
	return
}

func (m *matcher) MatchOther(line *Line) (ok bool, token *Token, err error) {
	token, ok = m.newTokenAtLocation(line.LineNumber, 0), true
	token.Type = TokenTypeOther

	element := line.LineText
	txt := strings.TrimLeft(element, " ")

	if len(element)-len(txt) > m.indentToRemove {
		token.Text = m.unescapeDocString(element[m.indentToRemove:])
	} else {
		token.Text = m.unescapeDocString(txt)
	}
	return
}

func (m *matcher) unescapeDocString(text string) string {
	if m.activeDocStringSeparator == DocstringSeparator {
		return strings.Replace(text, "\\\"\\\"\\\"", DocstringSeparator, -1)
	}
	if m.activeDocStringSeparator == DocstringAlternativeSeparator {
		return strings.Replace(text, "\\`\\`\\`", DocstringAlternativeSeparator, -1)
	}
	return text
}
