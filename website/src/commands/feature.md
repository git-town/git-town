# git town feature

<a type="command-summary">

```command-summary
git town feature [<branch-name>...] [-h | --help] [-v | --verbose]
```

</a>

The _feature_ command makes some of your branches
[feature](../branch-types.md#feature-branches) branches.

## Positional arguments

Make the current branch a feature branch:

```fish
git town feature
```

Make branches "alpha" and "beta" feature branches:

```fish
git town feature alpha beta
```

Check out a remote branch (that exists at the
[development remote](../preferences/dev-remote.md) but not on your local
machine) and make it a feature branch:

```fish
git town feature somebody-elses-branch
```

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [contribute](contribute.md) makes the chosen branches contribution branches
- [observe](observe.md) makes the chosen branches feature branches
- [prototype](prototype.md) makes the chosen branches prototype branches
- [park](park.md) makes the chosen branches parked
