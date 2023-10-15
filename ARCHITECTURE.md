# Git Town architecture

### Design goals

Git Town solves multiple large problems:

1. Execute a number of Git operations depending on conditions in the Git repo.
   Some of these conditions might change at runtime.
2. Allow the end user to resolve problems in the same terminal window and shell
   environment that Git Town executes in.
3. Reliably undo anything that Git Town has done upon request.

### Execution framework

Git Town addresses requirements 1 and 2 via an
[interpreter](https://en.wikipedia.org/wiki/Interpreter_(computing)) that
executes programs consisting of using Git-related opcodes. Each Git Town
command:

- inspects the state of the Git repo
- assembles a program that implements the Git operations that Git Town needs to
  perform
  - this program consists of opcodes that the Git Town interpreter can execute
- starts the Git Town interpreter engine to execute this program

If there are issues that require the user to resolve in a terminal window, the
interpreter:

- persists the current interpreter state (runstate) to disk
- exits the running Git Town process to lets the user use the terminal window
  and shell environment that they used to call Git Town to resolve the problems
- prints an explanation of the problem and what the user needs to do

After resolving the problems and restarting Git Town, the interpreter recognizes
and loads the persisted state from disk and resumes executing it.

### Undo framework

To undo a previously run Git Town command (requirement 3), Git Town:

- compares snapshots of the affected Git repository before and after the command
  ran
- determines the changes that the Git Town command made to the repo
- creates a program that reverses these changes
- starts the interpreter to execute this program
