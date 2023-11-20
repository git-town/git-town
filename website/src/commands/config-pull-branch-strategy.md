# git town config pull-branch-strategy <rebase|merge>

The _pull-branch-strategy_ configuration command displays or sets your pull
branch strategy. The pull branch strategy specifies which strategy to use when
merging remote tracking branches into local branches for the main branch and
perennial branches.

### Arguments

- without an argument, displays the current branch strategy
- with `rebase`, set the pull branch strategy to rebase
- with `merge`, set the pull branch strategy to merge
