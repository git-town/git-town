# Hosting platform

To talk to the API of your forge, Git Town needs to know which platform (GitHub,
Gitlab, Bitbucket, etc) you use.

By default, Git Town determines the forge by looking at the URL of the
[development remote](dev-remote.md). If that's not successful, for example when
using private instances of forges, you can tell Git Town through this
configuration setting which forge you use.

## values

You can use one of these values for the hosting platform setting:

- remove the entry or leave it empty for auto-detection
- `github`
- `gitlab`
- `gitea`
- `bitbucket`
- `bitbucket-datacenter`

## config file

In the [config file](../configuration-file.md) the hosting platform is part of
the `[hosting]` section:

```toml
[hosting]
platform = "<value>"
```

## Git metadata

To configure the hosting platform in Git, run this command:

```wrap
git config [--global] git-town.hosting-platform <value>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.
