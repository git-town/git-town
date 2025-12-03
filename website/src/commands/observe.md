# git town observe

<a type="git-town-command" />

```command-summary
git town observe [<branch-name>...] [-h | --help] [-v | --verbose]
```

The _observe_ command makes some of your branches
[observed](../branch-types.md#observed-branches) branches.

To convert an observed branch back into a feature branch, use the
[feature](feature.md) command.

## Positional arguments

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

## Options

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [contribute](contribute.md) makes the chosen branches contribution branches
- [feature](feature.md) makes the chosen branches feature branches
- [park](park.md) makes the chosen branches parked
- [prototype](prototype.md) makes the chosen branches prototype branches

<!-- keep-sorted end -->
