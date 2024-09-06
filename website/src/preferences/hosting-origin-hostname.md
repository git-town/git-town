# hosting.origin-hostname

If you use SSH identities, you can define the hostname of your source code
repository with this setting. The given value should match the hostname in your
SSH config file.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## config file

In the [config file](../configuration-file.md) the hosting platform is part of
the `[hosting]` section:

```toml
[hosting]
origin-hostname = "<hostname>"
```

## Git metadata

To configure the origin hostname in Git, run this command:

```bash
git config [--global] git-town.hosting-origin-hostname <hostname>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
