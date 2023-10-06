// Package subshell provides facilities to execute CLI commands in subshells.
package subshell

type Counter interface {
	RegisterRun()
}

type MessagesCollector interface{}
