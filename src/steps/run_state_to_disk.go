package steps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// LoadPreviousRunState loads the run state from disk if it exists or creates a new run state
func LoadPreviousRunState() (result *RunState, err error) {
	filename := getRunResultFilename()
	if util.DoesFileExist(filename) {
		var runState RunState
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return result, fmt.Errorf("cannot read file %q: %w", filename, err)
		}
		err = json.Unmarshal(content, &runState)
		if err != nil {
			return result, fmt.Errorf("cannot parse content of file %q: %w", filename, err)
		}
		return &runState, nil
	}
	return nil, nil
}

// DeletePreviousRunState deletes the previous run state from disk
func DeletePreviousRunState() error {
	filename := getRunResultFilename()
	if util.DoesFileExist(filename) {
		err := os.Remove(filename)
		if err != nil {
			return fmt.Errorf("cannot delete file %q: %w", filename, err)
		}
	}
	return nil
}

// SaveRunState saves the run state to disk
func SaveRunState(runState *RunState) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot encode run-state: %w", err)
	}
	filename := getRunResultFilename()
	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %w", filename, err)
	}
	return nil
}

func getRunResultFilename() string {
	replaceCharacterRegexp := regexp.MustCompile("[[:^alnum:]]")
	directory := replaceCharacterRegexp.ReplaceAllString(git.GetRootDirectory(), "-")
	return path.Join(os.TempDir(), directory)
}
