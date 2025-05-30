package runlog

type Event string

const (
	EventStart Event = "start" // a Git Town command has started
	EventEnd   Event = "end"   // a Git Town command has finished or errored
)
