# git town switch

The _switch_ command displays the branch hierarchy on your machine and allows
switching the current Git workspace to another local Git branch. Unlike
[git-switch](https://git-scm.com/docs/git-switch), Git Town's switch command
uses a more ergonomic visual UI and supports VIM motion commands.

`git town switch` does not allow switching to branches that are checked out in
other worktrees and notifies you about uncommitted changes in your workspace in
case you forgot to commit them to the current branch.

### Arguments

The `--merge` or `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag.
