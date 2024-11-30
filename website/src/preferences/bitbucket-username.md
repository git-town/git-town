# bitbucket-username

Git Town can interact with Bitbucket Cloud in your name, for example to update
pull requests as branches get created, shipped, or deleted. To do so, Git Town
needs your Bitbucket username and an
[Bitbucket App Password](bitbucket-app-password.md).

Bitbucket Datacenter capabilities are a little more restricted but the process
is the same as for Bitbucket Cloud.

The best way to enter your Bitbucket username is via the
[setup assistant](../configuration.md).

## config file

Since usernames are user specific, you cannot add them to the config file.

## Git metadata

You can configure the Bitbucket username manually by running:

```bash
git config [--global] git-town.bitbucket-username <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
