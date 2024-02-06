# hosting-platform

To talk to the API of your code hosting platform, Git Town needs to know which
platform (GitHub, Gitlab, Bitbucket, etc) you use.

Git Town can automatically figure out the code hosting platform by looking at
the URL of the `origin` remote. In cases where that's not successful, for
example when using private instances of code hosting platforms, you can tell Git
Town which code hosting platform you use.

The best way to configure this is via the
[setup assistant](../configuration.md).

## Configuration via config file

In the config file, the hosting platform is part of the `[hosting]` section:

```toml
[hosting]
platform = "<value>"
```

## Configuration via Git

To configure the hosting platform in Git:

```
git config [--global] git-town.hosting-platform <name>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.

## Values

The setting can have these values:

- remove the entry or leave it empty for auto-detection
- `github`
- `gitlab`
- `gitea`
- `bitbucket`
