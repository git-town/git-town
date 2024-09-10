# git town switch [--merge] [--all] [--type <branch type>] [regex...]

The _switch_ command displays the branch hierarchy on your machine and allows
switching the current Git workspace to another local Git branch. What
differentiates Git Town's switch command from
[git-switch](https://git-scm.com/docs/git-switch) is its ergonomic visual UI,
VIM motion commands, and filtering the list of branches to types and regular
expression matches.

`git town switch` notifies you about uncommitted changes in your workspace in
case you forgot to commit them to the current branch.

### arguments

`git town switch` interprets all positional arguments as regular expressions.
When receiving regular expressions from the user, it displays only the branches
that match at least one of the regular expressions.

### --merge

The `--merge` or `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag.

### --all

The `--all` or `-a` flag also displays both local and remote branches.

### --type

The `--type` or `-t` flag reduces the list of branches to those that have the
given type(s). For example, to display only observed branches:

Switch to one of your observed branches:

```
git town switch --type=observed
```

Branch types can be shortened:

```
git town switch -t o
```

This can be further compacted to:

```
git town switch -to
```

You can provide multiple branch types separated by `,`, `+`, `&`, or `|`, like
this:

```
git town switch --type=observed+contribution
```

This can be shortened to:

```
git town switch -to+c
```
