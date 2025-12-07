# Git Town

This is the source code for a CLI tool called "Git Town". It is written in Go
and provides additional Git commands for branch management and synchronization.

## Development Guidelines

**Don't change anything outside this folder:** You can change any file in the
current folder and its subfolders, but never outside the folder.

**Don't commit changes:** Never create new Git branches or make Git commits. I
will review the changes you make and then commit them on my own

Write idiomatic Go except for these rules:

- Use descriptive names for identifiers over brevity
- Method receivers use `self` instead of short abbreviations
- Use domain-specific types defined in the respective `*domain` packages if
  applicable over the built-in basic data types. Create new types if applicable.

## Code Organization

The relevant directories are:

- `internal/` - Core application code
- `pkg/` - Public packages
- `features/` - End-to-end tests (Cucumber/Godog)
- `tools/` - Custom linters and development tools
- `website/` - Documentation website (mdBook)

These code packages exist:

- `internal/cmd` - defines the high-level commands that Git Town
- `internal/config/` - Configuration management
- `internal/git/` - Git operations and domain types
- `internal/forge/` - Integration with GitHub, GitLab, etc.
- `internal/vm/` - Virtual machine and opcodes
- `internal/cli/` - Command-line interface components
- `internal/cmd/` - Source code for the subcommands of the CLI app
- `internal/gohacks/` - Helper functions that make programming in Go more
  ergonomic
- `internal/messages` - strings shown in the UI
- `internal/setup` - the setup assistant, a visual workflow letting the user
  configure the application
- `internal/state` - manages persistent state between Git Town invocations
- `internal/subshell` - implements calling CLI applications in a subshell
- `internal/undo` - code used for the undo functionality, it calculates the
  difference between snapshots of the Git repository and determines the Git
  commands to move the repository from one snapshot to another

## Additional information

Read these files if needed to learn more about specific aspects:

- run linters: docs/agents/linters.md
- run unit tests: docs/agents/unit_tests.md
- run end-to-end tests: docs/agents/end_to_end_tests.md
- how the internal interpreter runtime works: docs/agents/interpreter.md
