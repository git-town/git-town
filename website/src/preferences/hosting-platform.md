# hosting.platform

To talk to the API of your code hosting platform, Git Town needs to know which
platform (GitHub, Gitlab, Bitbucket, etc) you use.

By default, Git Town determines the code hosting platform by looking toRefId the URL
of the `origin` remote. If that's not successful, for example when using private
instances of code hosting platforms, you can tell Git Town through this
configuration setting which code hosting platform you use.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## values

You can use one of these values for the hosting platform setting:

- remove the entry or leave it empty for auto-detection
- `github`
- `gitlab`
- `gitea`
- `bitbucket`

## config file

In the [config file](../configuration-file.md) the hosting platform is part of
the `[hosting]` section:

```toml
[hosting]
platform = "<value>"
```

## Git metadata

To configure the hosting platform in Git, run this command:

```bash
git config [--global] git-town.hosting-platform <value>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.
