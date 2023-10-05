package runstate

import (
	"encoding/json"
	"fmt"
)

// CommandsCounter counts the commands run by a Git Town command.
type CommandsCounter struct {
	count int
}

// NewCommands provides new instances of Commands.
func NewCommands(count ...int) CommandsCounter {
	commands := CommandsCounter{count: 0}
	if len(count) > 0 {
		for i := 0; i < count[0]; i++ {
			commands.RegisterRun()
		}
	}
	return commands
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (c CommandsCounter) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.count)
}

func (c *CommandsCounter) RegisterRun() {
	c.count++
}

func (c *CommandsCounter) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", c.count)
}

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (c *CommandsCounter) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.count)
}
