# git town prototype

<a type="command-summary">

```command-summary
git town prototype [<branch-name>...] [-h | --help] [-v | --verbose]
```

</a>

The _prototype_ command marks some of your branches as
[prototype branches](../branch-types.md#prototype-branches).

To convert a prototype branch back into a feature branch, use the
[feature](feature.md) command.

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

<!-- keep-sorted start -->

- [contribute](contribute.md) makes the chosen branches contribution branches
- [feature](feature.md) makes the chosen branches feature branches
- [observe](observe.md) makes the chosen branches observed
- [park](park.md) makes the chosen branches parked

<!-- keep-sorted end -->
