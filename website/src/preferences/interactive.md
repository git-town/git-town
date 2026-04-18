# Interactive Mode

Git Town prompts for missing input when required CLI flags are not provided.

These prompts require an interactive terminal. If the terminal is
non-interactive or lacks required capabilities (for example, in CI or scripts),
Git Town disables interactive mode automatically.

If auto-detection is incorrect, or if you need consistent behavior, disable
interactive mode explicitly using this setting.

When interactive mode is disabled and required input is missing, Git Town exits
with an error message explaining the flags needed to rerun the command
non-interactively.

## via CLI flag

You can enable or disable interactive features for a single invocation:

```sh
git-town <command> --interactive
git-town <command> --non-interactive
```
