# git prototype [branches]

The _prototype_ command marks some of your branches as
[prototype branches](../branch-types.md#prototype-branches).

## Examples

Make the current branch a prototype branch:

```fish
git town prototype
```

Make branches "alpha" and "beta" prototype branches:

```fish
git town prototype alpha beta
```

Make the current prototype branch a feature branch:

```fish
git hack
```

The [hack](hack.md) command converts prototype branches back to feature
branches.
