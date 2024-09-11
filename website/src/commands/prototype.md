# git town prototype

> _git town prototype [branch-name...]_

The _prototype_ command marks some of your branches as
[prototype branches](../branch-types.md#prototype-branches).

To convert a contribution branch back into a feature branch, use the
[hack](hack.md) command.

### Positional arguments

Make the current branch a prototype branch:

```fish
git town prototype
```

Make branches "alpha" and "beta" prototype branches:

```fish
git town prototype alpha beta
```

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
