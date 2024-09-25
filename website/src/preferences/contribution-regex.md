# contribution-regex

Branches matching this regular expression are treated as
[contribution branches](../branch-types.md#contribution-branches).

## configure in config file

Defining the contribution regex in the [config file](../configuration-file.md)
is only a good idea if there the matching branches are considered contribution
for all team members. An example of that are branches created by an external
service like Renovate or Dependabot.

```toml
[branches]
contribution-regex = "^renovate/"
```

## configure in Git metadata

To manually set the feature regex, use the following command:

```bash
git config [--global] git-town.contribution-regex '^renovate/'
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
