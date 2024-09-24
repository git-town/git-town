# git town park

> _git town park [branch-name...]_

The _park_ command [parks](../branch-types.md#parked-branches) some of your
branches.

To convert a parked branch back into a feature branch, use the [hack](hack.md)
command.

## Positional arguments

Park the current branch:

```fish
git town park
```

Park branches "alpha" and "beta":

```fish
git town park alpha beta
```

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
