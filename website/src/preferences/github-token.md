# github-token

```
git-town.github-token=<token>
```

To interact with the GitHub API when [shipping](../commands/ship.md), Git Town
needs a
[personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the `repo` scope. You can create one in your
[account settings](https://github.com/settings/tokens/new). When you have
created the token, run the [setup assistant](../commands/config-setup.md) and
enter it there.

Alternatively, store the token manually by running
`git config git-town.github-token <token>` (where `<token>` is replaced with the
content of your GitHub access token) inside your code repository.
