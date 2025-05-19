package opcodes

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type ExitToShell struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}
