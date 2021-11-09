package formatters

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/cucumber/messages-go/v16"

	"github.com/cucumber/godog/colors"
	"github.com/cucumber/godog/internal/models"
	"github.com/cucumber/godog/internal/utils"
)

var (
	red    = colors.Red
	redb   = colors.Bold(colors.Red)
	green  = colors.Green
	blackb = colors.Bold(colors.Black)
	yellow = colors.Yellow
	cyan   = colors.Cyan
	cyanb  = colors.Bold(colors.Cyan)
	whiteb = colors.Bold(colors.White)
)

// repeats a space n times
var s = utils.S

var (
	passed    = models.Passed
	failed    = models.Failed
	skipped   = models.Skipped
	undefined = models.Undefined
	pending   = models.Pending
)

type sortFeaturesByName []*models.Feature

func (s sortFeaturesByName) Len() int           { return len(s) }
func (s sortFeaturesByName) Less(i, j int) bool { return s[i].Feature.Name < s[j].Feature.Name }
func (s sortFeaturesByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type sortPicklesByID []*messages.Pickle

func (s sortPicklesByID) Len() int { return len(s) }
func (s sortPicklesByID) Less(i, j int) bool {
	iID := mustConvertStringToInt(s[i].Id)
	jID := mustConvertStringToInt(s[j].Id)
	return iID < jID
}
func (s sortPicklesByID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type sortPickleStepResultsByPickleStepID []models.PickleStepResult

func (s sortPickleStepResultsByPickleStepID) Len() int { return len(s) }
func (s sortPickleStepResultsByPickleStepID) Less(i, j int) bool {
	iID := mustConvertStringToInt(s[i].PickleStepID)
	jID := mustConvertStringToInt(s[j].PickleStepID)
	return iID < jID
}
func (s sortPickleStepResultsByPickleStepID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func mustConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

// DefinitionID ...
func DefinitionID(sd *models.StepDefinition) string {
	ptr := sd.HandlerValue.Pointer()
	f := runtime.FuncForPC(ptr)
	file, line := f.FileLine(ptr)
	dir := filepath.Dir(file)

	fn := strings.Replace(f.Name(), dir, "", -1)
	var parts []string
	for _, gr := range matchFuncDefRef.FindAllStringSubmatch(fn, -1) {
		parts = append(parts, strings.Trim(gr[1], "_."))
	}
	if len(parts) > 0 {
		// case when suite is a structure with methods
		fn = strings.Join(parts, ".")
	} else {
		// case when steps are just plain funcs
		fn = strings.Trim(fn, "_.")
	}

	if pkg := os.Getenv("GODOG_TESTED_PACKAGE"); len(pkg) > 0 {
		fn = strings.Replace(fn, pkg, "", 1)
		fn = strings.TrimLeft(fn, ".")
		fn = strings.Replace(fn, "..", ".", -1)
	}

	return fmt.Sprintf("%s:%d -> %s", filepath.Base(file), line, fn)
}

var matchFuncDefRef = regexp.MustCompile(`\(([^\)]+)\)`)
