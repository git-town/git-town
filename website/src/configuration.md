# Set up configuration

If your repository already contains a `.git-town.toml` or `.git-branches.toml`
file, you're all set. If not - or if something isn't working as expected - Git
Town provides an interactive that guides you through the entire configuration
process. Just run:

```
git town config setup
```

This command walks you through all available configuration options, explains
what each one does, lets you adjust them, and validates that everything is
working correctly.

For more details on how Git Town handles configuration, see the
[configuration reference](preferences.md).

### API access

Some Git Town features require access the your code forge. This allows Git Town
to:

- infer the parent of a branch from open pull requests
- automatically update pull requests when you prepend, rename, or remove
  branches or change their parent
- trigger pull request merges directly from your terminal

Configuring API access is straightforward. Git Town supports the following
platforms:

- GitHub: Git Town can use GitHub's [gh CLI](https://cli.github.com) to talk to
  GitHub's API. If you don't have `gh` installed, you can also set up an
  [access token](preferences/github-token.md) and use Git Town's built-in GitHub
  connector.
- GitLab: Git Town can use GitLab's
  [glab CLI](https://gitlab.com/gitlab-org/cli/-/tree/main) to talk to GitLab's
  API. If you don't have this tool installed, you can also set up an
  [access token](preferences/gitlab-token.md) and use Git Town's built-in GitLab
  connector.
- Bitbucket: [username](preferences/bitbucket-username.md) and
  [app password](preferences/bitbucket-app-password.md)
- Gitea: [access token](preferences/gitea-token.md)
- Codeberg: [access token](preferences/codeberg-token.md)
