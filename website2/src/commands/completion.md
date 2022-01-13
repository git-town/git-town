# git town completions [bash|zsh|fish|powershell]

The _completions_ command generates auto-completion scripts for Bash, zsh, fish,
and PowerShell. With completions enabled, `git-town <tab key>` will show you all
possible subcommands, possibly with a short help text.

## Bash

To install autocompletion in Bash, run this command:

```
source <(git-town completions bash)
```

## Zsh

To install autocompletions in Zsh, run this command:

```
source <(git-town completions zsh)
```

## Fish

To install autocompletions in Fish, run this command:

```
git-town completions fish | source
```
