package statistics

import (
	"encoding/json"
	"fmt"
)

// Commands counts the commands run by a Git Town command.
type Commands struct {
	commandsRun int
	messages    []string
}

// NewCommands provides new instances of Commands for testing only!
func NewCommands(count int) Commands {
	commands := Commands{commandsRun: 0}
	for i := 0; i < count; i++ {
		commands.RegisterRun()
	}
	return commands
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (c Commands) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.commandsRun)
}

func (c *Commands) RegisterRun() {
	c.commandsRun++
}

func (c *Commands) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", c.commandsRun)
}

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (c *Commands) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.commandsRun)
}
