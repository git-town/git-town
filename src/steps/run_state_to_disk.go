package steps

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/pkg/errors"
)

// LoadPreviousRunState loads the run state from disk if it exists or creates a new run state
func LoadPreviousRunState() (result *RunState, err error) {
	filename := getRunResultFilename()
	if util.DoesFileExist(filename) {
		var runState RunState
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return result, errors.Wrapf(err, "cannot read file %q", filename)
		}
		err = json.Unmarshal(content, &runState)
		if err != nil {
			return result, errors.Wrapf(err, "cannot parse content of file %q", filename)
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
			return errors.Wrapf(err, "cannot delete file %q", filename)
		}
	}
	return nil
}

// SaveRunState saves the run state to disk
func SaveRunState(runState *RunState) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return errors.Wrap(err, "cannot encode run-state")
	}
	filename := getRunResultFilename()
	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return errors.Wrapf(err, "cannot write file %q", filename)
	}
	return nil
}

func getRunResultFilename() string {
	replaceCharacterRegexp := regexp.MustCompile("[[:^alnum:]]")
	directory := replaceCharacterRegexp.ReplaceAllString(git.GetRootDirectory(), "-")
	return path.Join(os.TempDir(), directory)
}
