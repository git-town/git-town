# contribution-regex

Branches matching this regular expression are treated as
[contribution branches](../branch-types.md#contribution-branches).

## configure in config file

In the [config file](../configuration-file.md), define the contribution regex
within the `[branches]` section. This is useful if external services create
proposals in your code base, which should be treated as contribution branches by
all team members.

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
