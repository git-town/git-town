package misspell

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func inArray(haystack []string, needle string) bool {
	for _, word := range haystack {
		if needle == word {
			return true
		}
	}
	return false
}

var wordRegexp = regexp.MustCompile(`[a-zA-Z0-9']+`)

// Diff is datastructure showing what changed in a single line
type Diff struct {
	Filename  string
	FullLine  string
	Line      int
	Column    int
	Original  string
	Corrected string
}

// Replacer is the main struct for spelling correction
type Replacer struct {
	Replacements []string
	Debug        bool
	engine       *strings.Replacer
	corrected    map[string]string
}

// New creates a new default Replacer using the main rule list
func New() *Replacer {
	r := Replacer{
		Replacements: DictMain,
	}
	r.Compile()
	return &r
}

// RemoveRule deletes existings rules.
// TODO: make inplace to save memory
func (r *Replacer) RemoveRule(ignore []string) {
	newwords := make([]string, 0, len(r.Replacements))
	for i := 0; i < len(r.Replacements); i += 2 {
		if inArray(ignore, r.Replacements[i]) {
			continue
		}
		newwords = append(newwords, r.Replacements[i:i+2]...)
	}
	r.engine = nil
	r.Replacements = newwords
}

// AddRuleList appends new rules.
// Input is in the same form as Strings.Replacer: [ old1, new1, old2, new2, ....]
// Note: does not check for duplictes
func (r *Replacer) AddRuleList(additions []string) {
	r.engine = nil
	r.Replacements = append(r.Replacements, additions...)
}

// Compile compiles the rules.  Required before using the Replace functions
func (r *Replacer) Compile() {

	r.corrected = make(map[string]string, len(r.Replacements)/2)
	for i := 0; i < len(r.Replacements); i += 2 {
		r.corrected[r.Replacements[i]] = r.Replacements[i+1]
	}
	r.engine = strings.NewReplacer(r.Replacements...)
}

/*
line1 and line2 are different
extract words from each line1

replace word -> newword
if word == new-word
  continue
if new-word in list of replacements
  continue
new word not original, and not in list of replacements
  some substring got mixed up.  UNdo
*/
func (r *Replacer) recheckLine(s string, lineNum int, buf io.Writer, next func(Diff)) {
	first := 0
	redacted := RemoveNotWords(s)

	idx := wordRegexp.FindAllStringIndex(redacted, -1)
	for _, ab := range idx {
		word := s[ab[0]:ab[1]]
		newword := r.engine.Replace(word)
		if newword == word {
			// no replacement done
			continue
		}
		if r.corrected[word] == newword {
			// word got corrected into something we know
			io.WriteString(buf, s[first:ab[0]])
			io.WriteString(buf, newword)
			first = ab[1]
			next(Diff{
				FullLine:  s,
				Line:      lineNum,
				Original:  word,
				Corrected: newword,
				Column:    ab[0],
			})
			continue
		}
		// Word got corrected into something unknown. Ignore it
	}
	io.WriteString(buf, s[first:])
}

// Replace is corrects misspellings in input, returning corrected version
//  along with a list of diffs.
func (r *Replacer) Replace(input string) (string, []Diff) {
	output := r.engine.Replace(input)
	if input == output {
		return input, nil
	}
	diffs := make([]Diff, 0, 8)
	buf := bytes.NewBuffer(make([]byte, 0, max(len(input), len(output))+100))
	// faster that making a bytes.Buffer and bufio.ReadString
	outlines := strings.SplitAfter(output, "\n")
	inlines := strings.SplitAfter(input, "\n")
	for i := 0; i < len(inlines); i++ {
		if inlines[i] == outlines[i] {
			buf.WriteString(outlines[i])
			continue
		}
		r.recheckLine(inlines[i], i+1, buf, func(d Diff) {
			diffs = append(diffs, d)
		})
	}

	return buf.String(), diffs
}

// ReplaceReader applies spelling corrections to a reader stream.  Diffs are
// emitted through a callback.
func (r *Replacer) ReplaceReader(raw io.Reader, w io.Writer, next func(Diff)) error {
	var (
		err     error
		line    string
		lineNum int
	)
	reader := bufio.NewReader(raw)
	for err == nil {
		lineNum++
		line, err = reader.ReadString('\n')

		// if it's EOF, then line has the last line
		// don't like the check of err here and
		// in for loop
		if err != nil && err != io.EOF {
			return err
		}
		// easily 5x faster than regexp+map
		if line == r.engine.Replace(line) {
			io.WriteString(w, line)
			continue
		}
		// but it can be inaccurate, so we need to double check
		r.recheckLine(line, lineNum, w, next)
	}
	return nil
}
