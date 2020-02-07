package godog

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/cucumber/godog/gherkin"
)

func init() {
	Format("junit", "Prints junit compatible xml to stdout", junitFunc)
}

func junitFunc(suite string, out io.Writer) Formatter {
	return &junitFormatter{
		basefmt: basefmt{
			suiteName: suite,
			started:   timeNowFunc(),
			indent:    2,
			out:       out,
		},
		lock: new(sync.Mutex),
	}
}

type junitFormatter struct {
	basefmt
	lock *sync.Mutex
}

func (f *junitFormatter) Node(n interface{}) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Node(n)
}

func (f *junitFormatter) Feature(ft *gherkin.Feature, p string, c []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Feature(ft, p, c)
}

func (f *junitFormatter) Summary() {
	suite := buildJUNITPackageSuite(f.suiteName, f.started, f.features)

	_, err := io.WriteString(f.out, xml.Header)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write junit string:", err)
	}

	enc := xml.NewEncoder(f.out)
	enc.Indent("", s(2))
	if err = enc.Encode(suite); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write junit xml:", err)
	}
}

func (f *junitFormatter) Passed(step *gherkin.Step, match *StepDef) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Passed(step, match)
}

func (f *junitFormatter) Skipped(step *gherkin.Step, match *StepDef) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Skipped(step, match)
}

func (f *junitFormatter) Undefined(step *gherkin.Step, match *StepDef) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Undefined(step, match)
}

func (f *junitFormatter) Failed(step *gherkin.Step, match *StepDef, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Failed(step, match, err)
}

func (f *junitFormatter) Pending(step *gherkin.Step, match *StepDef) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.basefmt.Pending(step, match)
}

func (f *junitFormatter) Sync(cf ConcurrentFormatter) {
	if source, ok := cf.(*junitFormatter); ok {
		f.lock = source.lock
	}
}

func (f *junitFormatter) Copy(cf ConcurrentFormatter) {
	if source, ok := cf.(*junitFormatter); ok {
		for _, v := range source.features {
			f.features = append(f.features, v)
		}
		for _, v := range source.failed {
			f.failed = append(f.failed, v)
		}
		for _, v := range source.passed {
			f.passed = append(f.passed, v)
		}
		for _, v := range source.skipped {
			f.skipped = append(f.skipped, v)
		}
		for _, v := range source.undefined {
			f.undefined = append(f.undefined, v)
		}
		for _, v := range source.pending {
			f.pending = append(f.pending, v)
		}
	}
}

func buildJUNITPackageSuite(suiteName string, startedAt time.Time, features []*feature) junitPackageSuite {
	suite := junitPackageSuite{
		Name:       suiteName,
		TestSuites: make([]*junitTestSuite, len(features)),
		Time:       timeNowFunc().Sub(startedAt).String(),
	}

	sort.Sort(sortByName(features))

	for idx, feat := range features {
		ts := junitTestSuite{
			Name:      feat.Name,
			Time:      feat.finishedAt().Sub(feat.startedAt()).String(),
			TestCases: make([]*junitTestCase, len(feat.Scenarios)),
		}

		for idx, scenario := range feat.Scenarios {
			tc := junitTestCase{}
			tc.Name = scenario.Name
			tc.Time = scenario.finishedAt().Sub(scenario.startedAt()).String()

			ts.Tests++
			suite.Tests++

			for _, step := range scenario.Steps {
				switch step.typ {
				case passed:
					tc.Status = passed.String()
				case failed:
					tc.Status = failed.String()
					tc.Failure = &junitFailure{
						Message: fmt.Sprintf("%s %s: %s", step.step.Type, step.step.Text, step.err),
					}
				case skipped:
					tc.Error = append(tc.Error, &junitError{
						Type:    "skipped",
						Message: fmt.Sprintf("%s %s", step.step.Type, step.step.Text),
					})
				case undefined:
					tc.Status = undefined.String()
					tc.Error = append(tc.Error, &junitError{
						Type:    "undefined",
						Message: fmt.Sprintf("%s %s", step.step.Type, step.step.Text),
					})
				case pending:
					tc.Status = pending.String()
					tc.Error = append(tc.Error, &junitError{
						Type:    "pending",
						Message: fmt.Sprintf("%s %s: TODO: write pending definition", step.step.Type, step.step.Text),
					})
				}
			}

			switch tc.Status {
			case failed.String():
				ts.Failures++
				suite.Failures++
			case undefined.String(), pending.String():
				ts.Errors++
				suite.Errors++
			}

			ts.TestCases[idx] = &tc
		}

		suite.TestSuites[idx] = &ts
	}

	return suite
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
}

type junitError struct {
	XMLName xml.Name `xml:"error,omitempty"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
}

type junitTestCase struct {
	XMLName xml.Name      `xml:"testcase"`
	Name    string        `xml:"name,attr"`
	Status  string        `xml:"status,attr"`
	Time    string        `xml:"time,attr"`
	Failure *junitFailure `xml:"failure,omitempty"`
	Error   []*junitError
}

type junitTestSuite struct {
	XMLName   xml.Name `xml:"testsuite"`
	Name      string   `xml:"name,attr"`
	Tests     int      `xml:"tests,attr"`
	Skipped   int      `xml:"skipped,attr"`
	Failures  int      `xml:"failures,attr"`
	Errors    int      `xml:"errors,attr"`
	Time      string   `xml:"time,attr"`
	TestCases []*junitTestCase
}

type junitPackageSuite struct {
	XMLName    xml.Name `xml:"testsuites"`
	Name       string   `xml:"name,attr"`
	Tests      int      `xml:"tests,attr"`
	Skipped    int      `xml:"skipped,attr"`
	Failures   int      `xml:"failures,attr"`
	Errors     int      `xml:"errors,attr"`
	Time       string   `xml:"time,attr"`
	TestSuites []*junitTestSuite
}
