# Interactive

By default, when Git Town needs additional information from the user that wasn't
provided via CLI flags, it asks directly for them via interactive dialogs. These
interactive dialogs don't work in environments that have limited or no
interactive terminal available. Git Town automatically disables interactivity
when it detects missing or degraded terminal features. If that automated
detection doesn't work in your situation, or you want to always disable
interactivity, this setting is for you.

When interactivity is disabled, and Git Town needs additional information, it
exits with an error message that describes the CLI flags to call Git Town with
again to provide the needed information.

## via CLI flag

You can enable or disable interactive features for a single invocation:

```sh
git-town <command> --interactive
git-town <command> --non-interactive
```
