package steps

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// LoadPreviousRunState loads the run state from disk if it exists or creates a new run state
func LoadPreviousRunState() *RunState {
	filename := getRunResultFilename()
	if util.DoesFileExist(filename) {
		var runState RunState
		content, err := ioutil.ReadFile(filename)
		exit.If(err)
		err = json.Unmarshal(content, &runState)
		exit.If(err)
		return &runState
	}
	return nil
}

// DeletePreviousRunState deletes the previous run state from disk
func DeletePreviousRunState() {
	filename := getRunResultFilename()
	if util.DoesFileExist(filename) {
		exit.If(os.Remove(filename))
	}
}

// SaveRunState saves the run state to disk
func SaveRunState(runState *RunState) {
	content, err := json.MarshalIndent(runState, "", "  ")
	exit.If(err)
	filename := getRunResultFilename()
	err = ioutil.WriteFile(filename, content, 0644)
	exit.If(err)
}

func getRunResultFilename() string {
	replaceCharacterRegexp, err := regexp.Compile("[[:^alnum:]]")
	exit.IfWrap(err, "Error compiling replace character expression")
	directory := replaceCharacterRegexp.ReplaceAllString(git.GetRootDirectory(), "-")
	return path.Join(os.TempDir(), directory)
}
