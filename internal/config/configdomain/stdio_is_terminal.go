package configdomain

// StdioIsTerminal indicates whether both STDIN and STDOUT of this process are connected to a terminal.
type StdioIsTerminal bool

func (self StdioIsTerminal) IsTerminal() bool {
	return bool(self)
}
