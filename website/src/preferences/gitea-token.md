# gitea-token

Git Town can interact with Gitea in your name, for example to update pull
requests as branches get created, shipped, or deleted. To do so, Git Town needs
a personal access token for Gitea.

The best way to enter your token is via the
[setup assistant](../configuration.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the API token manually by running:

```bash
git config [--global] git-town.gitea-token <token>
```
