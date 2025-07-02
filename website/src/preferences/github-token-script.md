# GitHub token

Git Town can interact with GitHub in your name, for example to update pull
requests as branches get created, shipped, or deleted. To do so, Git Town needs
a
[personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the `repo` scope. You can create one in your
[account settings](https://github.com/settings/tokens/new) or get one created
for you by using the [`gh` connector type](github-connector.md).

If you store such tokens securely, for example using a password manager, you can
provide a command that Git Town executes to obtain the GitHub token.

## config file

```toml
[hosting]
github-token-script = "<your script here>"
```

## Git metadata

You can configure the API token manually by running:

```wrap
git config [--global] git-town.github-token-script <script>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
