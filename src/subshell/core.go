// Package subshell provides facilities to execute CLI commands in subshells.
package subshell

type CommandsStats interface {
	RegisterRun()
}

type MessagesCollector interface {
}
