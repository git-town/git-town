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
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	tagTime, err := time.Parse(time.RFC3339, strings.TrimSpace(string(output)))
	if err != nil {
		panic(err.Error())
	}
	return Tag{
		ISOTime: tagTime.Format("2006-01-02"), // the time this tag was created, in ISO format
		Name:    name,                         // name of the tag
		Time:    tagTime,                      // the time this tag was created
	}
}
