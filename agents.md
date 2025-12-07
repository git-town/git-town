# Git Town

This is the source code for a CLI tool called "Git Town". It is written in Go
and provides additional Git commands for branch management and synchronization.

## Development Guidelines

You can change any file in the current folder and its subfolders, but never
outside the folder.

Never create new Git branches or make Git commits. I will review the changes you
make and then commit them on my own.

## Code Organization

The codebase is organized into orthogonal subsystems with `*domain` packages:

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

#### Important Directories

- `internal/` - Core application code
- `pkg/` - Public packages
- `features/` - End-to-end tests (Cucumber/Godog)
- `tools/` - Custom linters and development tools
- `website/` - Documentation website (mdBook)

### Code Style

Write idiomatic Go with these exceptions:

- Use descriptive names for identifiers over brevity
- Method receivers use `self` instead of short abbreviations
- Use domain-specific types defined in the respective `*domain` packages if
  applicable over the built-in basic type. Create new types if applicable.
