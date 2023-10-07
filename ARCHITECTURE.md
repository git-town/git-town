# Git Town architecture

### Overview

Git Town solves two large problems:

1. Execute a number of Git operations depending on conditions in the Git repo.
   These conditions might change at runtime.
2. Pause and resume this execution across several invocations of Git Town to let
   the user resolve merge conflicts in the same terminal window and shell
   environment that Git Town runs in.
3. Reliably undo anything that Git Town has done at the user's request.

### Execution framework

Git Town addresses requirements 1 and 2 via an
[interpreter](https://en.wikipedia.org/wiki/Interpreter_(computing)) that
executes programs written using Git-related operations. Each Git Town command:

- analyzes the state of the Git repo
- creates a "program" that implements the Git operations that Git Town needs to
  perform
- starts the interpreter to execute this program

If there are issues that require the user to resolve in a terminal window, the
interpreter:

- persists the current state to disk
- exits the running Git Town process to lets the user use the terminal window
  and shell environment that Git Town was running in to resolve the problems
- prints an explanation of the problem and what the user needs to do

After resolving the problems and restarting Git Town, the interpreter loads the
persisted state from disk and resumes execution.

### Undo framework

To undo a previously run Git Town command, Git Town:

- compares snapshots of the affected Git repository before and after the command
  ran
- determines the changes that the previously running Git Town command made to
  the repo
- create a program that reverses these changes
- starts the interpreter to execute this program
