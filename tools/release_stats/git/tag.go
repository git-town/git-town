package git

import (
	"os/exec"
	"strings"
	"time"
)

// all the information we need to know about Git tags in the context of this program
type Tag struct {
	ISOTime string
	Name    string
	Time    time.Time
}

// provides the time when the Git tag with the given name was created
func LoadTag(name string) Tag {
	cmd := exec.Command("git", "log", "-1", "--format=%cI", name)
	outputData, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	output := strings.TrimSpace(string(outputData))
	gitTime, err := time.Parse(time.RFC3339, output)
	if err != nil {
		panic(err.Error())
	}
	return Tag{
		ISOTime: gitTime.Format("2006-01-02"),
		Name:    name,
		Time:    gitTime,
	}
}
