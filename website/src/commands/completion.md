# git town completion [bash|zsh|fish|powershell]

The _completions_ command generates auto-completion scripts for Bash, zsh, fish,
and PowerShell. With completions enabled, `git-town <tab key>` will show you all
possible subcommands, possibly with a short help text.

## Bash

To install autocompletion for Bash, run this command:

```
source <(git-town completion bash)
```

## Zsh

To install autocompletions for Zsh, run this command:

```
source <(git-town completion zsh)
```

## Fish

To install autocompletions for Fish, run this command:

```
git-town completion fish | source
```

## Powershell

To install autocompletions for Powershell, run this command:

```
git-town completion powershell | source
```
