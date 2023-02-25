# git switch

The _switch_ command allows switching the current Git workspace to another local
Git branch. Unlike [git-switch](https://git-scm.com/docs/git-switch), Git Town's
switch command uses a more ergonomic visual UI. You can use these keys to
navigate the UI:

- `UP`, `k`: move the selection up
- `DOWN`, `TAB`, `j`: move the selection down
- `ENTER`, `s`: switch to the selected branch
- `ESC`: abort the dialog

The dialog starts at the selected branch, so running
`git town switch <ENTER><ENTER>` displays the branch hierarchy, with the active
branch highlighted.
