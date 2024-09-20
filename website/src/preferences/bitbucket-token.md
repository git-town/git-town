# bitbucket-token

Git Town can interact with Bitbucket Cloud in your name, for example to update
pull requests as branches get created, shipped, or deleted. To do so, Git Town
needs an app password for BitBucket.

To create an app password, click on the `Settings` cogwheel, choose
`Personal Bitbucket settings`, and then in the menu on the left `App passwords`.
You need an app password with these permissions:

- pull requests: read and write

The best way to enter your token is via the
[setup assistant](../configuration.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the API token manually by running:

```bash
git config [--global] git-town.gitea-token <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
