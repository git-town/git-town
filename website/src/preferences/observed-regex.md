# observed-regex

Branches matching this regular expression are treated as
[observed branches](../branch-types.md#observed-branches).

## configure in config file

Defining the observed regex in the [config file](../configuration-file.md) is
only a good idea if there the matching branches are considered observed for all
team members. An example of that are branches created by an external service
like Renovate or Dependabot.

```toml
[branches]
observed-regex = "^renovate/"
```

## configure in Git metadata

To manually set the feature regex, use the following command:

```bash
git config [--global] git-town.observed-regex '^renovate/'
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
