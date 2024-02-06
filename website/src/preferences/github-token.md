# github-token

To interact with the GitHub API, Git Town needs a
[personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the `repo` scope. You can create one in your
[account settings](https://github.com/settings/tokens/new).

The best way to enter your token is via the
[setup assistant](../configuration.md). Since your API token is confidential,
you cannot enter it into the config file.

## in Git metadata

You can configure the API token manually by running:

```bash
git config [--global] git-town.github-token <token>
```
