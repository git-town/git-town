package formatters

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/cucumber/messages-go/v16"
)

// some snippet formatting regexps
var snippetExprCleanup = regexp.MustCompile("([\\/\\[\\]\\(\\)\\\\^\\$\\.\\|\\?\\*\\+\\'])")
var snippetExprQuoted = regexp.MustCompile("(\\W|^)\"(?:[^\"]*)\"(\\W|$)")
var snippetMethodName = regexp.MustCompile("[^a-zA-Z\\_\\ ]")
var snippetNumbers = regexp.MustCompile("(\\d+)")

var snippetHelperFuncs = template.FuncMap{
	"backticked": func(s string) string {
		return "`" + s + "`"
	},
}

var undefinedSnippetsTpl = template.Must(template.New("snippets").Funcs(snippetHelperFuncs).Parse(`
{{ range . }}func {{ .Method }}({{ .Args }}) error {
	return godog.ErrPending
}

{{end}}func InitializeScenario(ctx *godog.ScenarioContext) { {{ range . }}
	ctx.Step({{ backticked .Expr }}, {{ .Method }}){{end}}
}
`))

type undefinedSnippet struct {
	Method   string
	Expr     string
	argument *messages.PickleStepArgument
}

func (s undefinedSnippet) Args() (ret string) {
	var (
		args      []string
		pos       int
		breakLoop bool
	)

	for !breakLoop {
		part := s.Expr[pos:]
		ipos := strings.Index(part, "(\\d+)")
		spos := strings.Index(part, "\"([^\"]*)\"")

		switch {
		case spos == -1 && ipos == -1:
			breakLoop = true
		case spos == -1:
			pos += ipos + len("(\\d+)")
			args = append(args, reflect.Int.String())
		case ipos == -1:
			pos += spos + len("\"([^\"]*)\"")
			args = append(args, reflect.String.String())
		case ipos < spos:
			pos += ipos + len("(\\d+)")
			args = append(args, reflect.Int.String())
		case spos < ipos:
			pos += spos + len("\"([^\"]*)\"")
			args = append(args, reflect.String.String())
		}
	}

	if s.argument != nil {
		if s.argument.DocString != nil {
			args = append(args, "*godog.DocString")
		}

		if s.argument.DataTable != nil {
			args = append(args, "*godog.Table")
		}
	}

	var last string

	for i, arg := range args {
		if last == "" || last == arg {
			ret += fmt.Sprintf("arg%d, ", i+1)
		} else {
			ret = strings.TrimRight(ret, ", ") + fmt.Sprintf(" %s, arg%d, ", last, i+1)
		}

		last = arg
	}

	return strings.TrimSpace(strings.TrimRight(ret, ", ") + " " + last)
}

type snippetSortByMethod []undefinedSnippet

func (s snippetSortByMethod) Len() int {
	return len(s)
}

func (s snippetSortByMethod) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s snippetSortByMethod) Less(i, j int) bool {
	return s[i].Method < s[j].Method
}
