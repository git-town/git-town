# git town switch

<a type="command-summary">

```command-summary
git town switch [<branch-name-regex>...] [-a | --all] [-d | --display-types <type>] [-h | --help] [-m | --merge] [-o | --order <asc|desc>] [-t | --type <name>] [-v | --verbose]
```

</a>

The _switch_ command displays the branch hierarchy on your machine and allows
switching the current Git workspace to another local Git branch using VIM motion
commands. It can filter the list of branches to particular branch types and
regular expression matches.

`git town switch` reminds you about uncommitted changes in your workspace in
case you forgot to commit them to the current branch.

## Positional arguments

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

## Options

#### `-a`<br>`--all`

The `--all` aka `-a` flag also displays both local and remote branches.

#### `-d`<br>`--display-types`

This flag allows customizing whether Git Town also displays the branch type in
addition to the branch name when showing a list of branches. More info
[here](../preferences/display-types.md#cli-flags).

#### `-m`<br>`--merge`

The `--merge` aka `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag. It attempts to merge uncommitted changes in your workspace into the target
branch.

This is useful when you have uncommitted changes in your current branch and want
to move them to the new branch.

#### `-o`<br>`--order`

The `--order` flag allows customizing the order in which branches get displayed.
More info [here](../preferences/order.md#cli-flag)

#### `-t <name>`<br>`--type <name>`

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

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [branch](branch.md) displays the branch hierarchy
- [walk](walk.md) executes a shell command or opens a shell in each of your
  local branches
