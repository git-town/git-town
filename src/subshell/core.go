// Package subshell provides facilities to execute CLI commands in subshells.
package subshell

type Statistics interface {
	RegisterRun()
}
