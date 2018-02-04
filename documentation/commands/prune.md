#### NAME

prune - clean up outdated information

#### SYNOPSIS

```
git town prune [branches | config]
git town prune (--abort | --continue)
```

#### DESCRIPTION

`git-town prune branches` deletes branches whose tracking branch no longer exists
from the local repository.
This usually means the branch was shipped or killed on another machine.

`git-town prune config` subcommand deletes branches from the local Git configuration
that don't exist in the local workspace.

`git-town prune` runs all subcommands.


#### OPTIONS

```
--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```

#### SEE ALSO
* [git config --reset](config.md) to remove all Git Town configuration from the current repository
