# git town contribute

<a type="command-summary">

```command-summary
git town contribute [<branch-name>...] [-h | --help] [-v | --verbose]
```

</a>

The _contribute_ command makes some of your branches
[contribution](../branch-types.md#contribution-branches) branches.

When called without arguments, it makes the current branch a contribution
branch.

To convert a contribution branch back into a feature branch, use the
[feature](feature.md) command.

To make the current branch a contribution branch:

```fish
git town contribute
```

## Positional arguments

When called with positional arguments, this commands makes the branches with the
given names contribution branches.

To make branches "alpha" and "beta" contribution branches:

```fish
git town contribute alpha beta
```

Check out a remote branch (that exists at the
[development remote](../preferences/dev-remote.md) but not on your local
machine) and make it a contribution branch:

```fish
git town contribute somebody-elses-branch
```

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [feature](feature.md) makes the chosen branches feature branches
- [observe](observe.md) makes the chosen branches observed
- [park](park.md) makes the chosen branches parked
- [prototype](prototype.md) makes the chosen branches prototype branches

<!-- keep-sorted end -->
