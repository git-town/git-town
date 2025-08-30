# Origin hostname

If you use SSH identities, you can define the hostname of your source code
repository with this setting. The given value should match the hostname in your
SSH config file.

## config file

In the [config file](../configuration-file.md) the forge is part of the
`[hosting]` section:

```toml
[hosting]
origin-hostname = "<hostname>"
```

## Git metadata

To configure the origin hostname in Git, run this command:

```wrap
git config [--global] git-town.hosting-origin-hostname <hostname>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the origin hostname by setting the `GIT_TOWN_ORIGIN_HOSTNAME`
environment variable.
