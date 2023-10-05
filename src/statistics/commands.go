package statistics

import "fmt"

// Commands is a Statistics implementation that counts how many commands were run.
type Commands struct {
	CommandsCount int
}

func (c *Commands) RegisterRun() {
	c.CommandsCount++
}

func (c *Commands) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", c.CommandsCount)
}
