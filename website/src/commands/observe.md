# git town observe

> _git town observe [branch-name...]_

The _observe_ command makes some of your branches
[observed](../branch-types.md#observed-branches) branches.

To convert an observed branch back into a feature branch, use the
[hack](hack.md) command.

### Positional arguments

Observe the current branch:

```fish
git town observe
```

Observe branches "alpha" and "beta":

```fish
git town observe alpha beta
```

Check out a remote branch (that exists at the
[development remote](../preferences/dev-remote.md) but not on your local
machine) and make it observed:

```fish
git town observe somebody-elses-branch
```

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
