# git town config [subcommand]

The _config_ command displays and updates the local Git Town configuration.

### Subcommands

Running without a subcommand shows the current Git Town configuration.

- The [get-parent](config-get-parent.md) subcommand prints the parent branch of
  the current or given branch.
- The [remove](config-remove.md) subcommand deletes all Git Town configuration
  entries.
- The [setup](config-setup.md) subcommand interactively prompts for all
  configuration values

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
