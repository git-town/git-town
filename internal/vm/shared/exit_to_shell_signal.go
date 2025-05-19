package shared

// ExitToShellError is a special error type that signals that no error happened
// and Git Town should simply exit to the shell without an error code,
// allowing resume via "git town continue".
type ExitToShellError struct{}

func (self ExitToShellError) Error() string {
	return ""
}

var ErrExitToShell = ExitToShellError{}
