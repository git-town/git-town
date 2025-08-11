# git town prototype

```command-summary
git town prototype [<branch-name>...] [-v | --verbose]
```

The _prototype_ command marks some of your branches as
[prototype branches](../branch-types.md#prototype-branches).

To convert a prototype branch back into a feature branch, use the
[hack](hack.md) command.

## Positional arguments

Make the current branch a prototype branch:

```fish
git town prototype
```

Make branches "alpha" and "beta" prototype branches:

```fish
git town prototype alpha beta
```

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [contribute](contribute.md) makes the chosen branches contribution branches
- [hack](hack.md) makes the chosen branches feature branches
- [observe](observe.md) makes the chosen branches observed
- [park](park.md) makes the chosen branches parked
