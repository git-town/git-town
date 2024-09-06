# feature-regex

Branches matching this regular expression are treated as feature branches. This
setting is relevant only when the [default-branch-type](default-branch-type.md)
setting is set to something different than "feature".

## configure in config file

In the [config file](../configuration-file.md), define the feature regex within
the `[branches]` section:

```toml
[branches]
feature-regex = "^my-*"
```

## configure in Git metadata

To manually set the feature regex, use the following command:

```bash
git config [--global] git-town.feature-regex '^user-.*'
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
