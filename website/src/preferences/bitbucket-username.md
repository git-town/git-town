# bitbucket-username

Git Town can interact with Bitbucket Cloud in your name, for example to update
pull requests as branches get created, shipped, or deleted. To do so, Git Town
needs your Bitbucket username and an
[Bitbucket App Password](bitbucket-app-password.md).

The best way to enter your Bitbucket username is via the
[setup assistant](../configuration.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the app password manually by running:

```bash
git config [--global] git-town.bitbucket-username <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
