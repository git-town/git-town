#### NAME

prune - clean up outdated information

#### SYNOPSIS

```
git town prune [<branches|config>]
```

#### DESCRIPTION

The "branches" subcommand deletes branches whose tracking branch no longer exists
from the local repository.
This usually means the branch was shipped or killed on another machine.

The "config" subcommand deletes branches from the local Git configuration
that don't exist in the local workspace.

Running "git-town prune" runs all subcommands.
