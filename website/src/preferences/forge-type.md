# Forge Type

To talk to the API of your forge, Git Town needs to know which forge type
(GitHub, Gitlab, Bitbucket, gitea, Forgejo, etc) you use.

By default, Git Town determines the forge type by looking at the URL of the
[development remote](dev-remote.md). If that's not successful, for example when
using a private forge, you can tell Git Town through this configuration setting
which forge type you use.

## values

You can use one of these values for the forge type setting:

- remove the entry or leave it empty for auto-detection
- `github`
- `gitlab`
- `gitea`
- `bitbucket`
- `bitbucket-datacenter`
- `forgejo`
- `azuredevops` (experimental)

## config file

In the [config file](../configuration-file.md) the forge type is part of the
`[hosting]` section:

```toml
[hosting]
forge-type = "<value>"
```

## Git metadata

To configure the forge type in Git, run this command:

```wrap
git config [--global] git-town.forge-type <value>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.

## environment variable

You can configure the forge type by setting the `GIT_TOWN_FORGE_TYPE`
environment variable.
