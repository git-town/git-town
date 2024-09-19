# gitea-token

Git Town can interact with Gitea in your name, for example to update pull
requests as branches get created, shipped, or deleted. To do so, Git Town needs
a personal access token for Gitea.

To create an API token, click on your profile image, choose `Settings`, and then
in the menu on the left `Applications`. You need an API token with permissions
to read the repository and issues.

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
