// Package subshell provides facilities to execute CLI commands in subshells.
//
// There are are two types of shell commands in Git Town:
//
//  1. Internal shell commands.
//     Git Town runs these silently to determine the state of a Git repository.
//     Git Town needs to know the output that they generated.
//     These commands don't change the Git repository, they only investigate it.
//
//  2. Public shell commands.
//     These are the commands that Git Town runs for the end user to change their Git repository.
//     Git Town doesn't need to know their output, only whether they failed.
//
// Package subshell provides various facilities to run internal and public shell commands.
package subshell
