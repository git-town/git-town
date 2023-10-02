package statistics

import "fmt"

// CommandsStatistics is a Statistics implementation that counts how many commands were run.
type CommandsStatistics struct {
	CommandsCount int
}

func (cs *CommandsStatistics) RegisterRun() {
	cs.CommandsCount++
}

func (cs *CommandsStatistics) PrintAnalysis() {
	fmt.Printf("Ran %d shell commands.", cs.CommandsCount)
}
