# git town park

<a type="command-summary">

```command-summary
git town park [<branch-name>...] [-h | --help] [-v | --verbose]
```

</a>

The _park_ command [parks](../branch-types.md#parked-branches) some of your
branches.

To convert a parked branch back into a feature branch, use the
[feature](feature.md) command or [propose](propose.md) it.

## Positional arguments

Park the current branch:

```fish
git town park
```

Park branches "alpha" and "beta":

```fish
git town park alpha beta
```

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [contribute](contribute.md) makes the chosen branches contribution branches
- [feature](feature.md) makes the chosen branches feature branches
- [observe](observe.md) makes the chosen branches observed
- [prototype](prototype.md) makes the chosen branches prototype branches
