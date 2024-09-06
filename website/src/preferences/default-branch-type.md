# default-branch-type

This setting defines the default branch type for Git Town when the branch type
is unknown. It applies to branches not listed in [main-branch](main-branch.md),
[perennial-branches](perennial-branches.md),
[contribution-branches](contribution-branches.md),
[observed-branches](observed-branches.md),
[parked-branches](parked-branches.md), or
[prototype-branches](prototype-branches.md).

Possible values are:

- `feature` (default)
- `contribution`
- `observed`
- `parked`
- `prototype`

If you set this to anything other than `feature`, you must also configure the
[feature-regex](feature-regex.md) setting. Otherwise, there will be no feature
branches.

## configuration via setup assistant

A great way to configure this setting is through the setup assistant.

## configure in config file

In the [config file](../configuration-file.md), the default branch type is
specified in the `[branches]` section:

```toml
[branches]
default-type = "feature"
```

## configure in Git metadata

You can manually configure the default branch type using Git metadata:

```bash
git config [--global] git-town.default-branch-type "feature"
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
