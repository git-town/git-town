# git town switch

> _git town switch [--merge] [--all] [--type] [branch-name-regex...]_

The _switch_ command displays the branch hierarchy on your machine and allows
switching the current Git workspace to another local Git branch using VIM motion
commands. It can filter the list of branches to particular branch types and
regular expression matches.

`git town switch` reminds you about uncommitted changes in your workspace in
case you forgot to commit them to the current branch.

### positional arguments

`git town switch` interprets all positional arguments as regular expressions.
When receiving regular expressions from the user, it displays only the branches
that match at least one of the regular expressions.

As an example, assuming all your branches start with `me-`, you can use this
command to switch to one of them:

```
git town switch me-
```

To display all branches starting with `me-` and the main branch:

```
git town switch me- main
```

### --all / -a

The `--all` aka `-a` flag also displays both local and remote branches.

### --display-types / -d

When enabled, this command displays the types for all branches except the main
branch and feature branches.

### --merge / -m

The `--merge` aka `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag.

### --type / -t &lt;name&gt;

The `--type` aka `-t` flag reduces the list of branches to those that have the
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

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
