package runlog

type Event string

const (
	EventStart Event = "start" // before a Git Town command starts making changes to the repo
	EventEnd   Event = "end"   // after a Git Town command finished making changes to the repo
)
