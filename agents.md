# Git Town

Git Town is a CLI tool written in Go that provides additional Git commands for
automating branch management, synchronization, and cleanup.

## Development Guidelines

- you can change any file in the current folder and its subfolders
- never change files outside the Git repository
- never create new Git branches
- never make Git commits
- I will review the changes you make and then commit them on my own.

## Automated testing

To run all unit tests for the project, use this command:

```bash
make unit
```

To run a single Cucumber test, also called end-to-end test, use this command:

```bash
make install
go test -- <test path>
```

### Linters

Please execute the linters after making changes to verify the correctness of
your changes. To run them, use the following command:

```bash
make lint
```

### End-to-End Tests

End-to-end tests are defined in the "features" directory. They take a while to
execute, so only run them to verify that everything still works after you are
done making changes. To run all end-to-end tests:

```bash
make cuke
```

## Key Architectural Components

#### VM-Based Execution Framework

Git Town uses an interpreter that executes self-modifying code consisting of
Git-related opcodes:

- Commands inspect Git repo state and generate a program of opcodes
- The interpreter (`internal/vm/`) executes these programs
- Programs can modify themselves at runtime based on repo state
- Runstate is persisted to disk for resume capability

#### Subsystem Organization

The codebase is organized into orthogonal subsystems with `*domain` packages:

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

- Write idiomatic Go except for the conditions listed below
- Use descriptive naming over brevity
- Method receivers use `self` instead of short abbreviations
- Use domain-specific types defined in the respective `*domain` packages if
  applicable over the built-in basic type.

## Common Development Tasks

### Debugging End-to-End Tests

- Add the `@this` tag to a specific scenario and then run `make cukethis` to
  execute only the tagged scenario
- Add the `@debug` tag to a Cucumber scenario to see CLI output and Git commands
