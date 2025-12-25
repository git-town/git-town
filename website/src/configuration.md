# Set up configuration

If your repository already contains a `git-town.toml`, `.git-town.toml`, or
`.git-branches.toml` file, you're all set. If not - or if something isn't
working as expected - Git Town provides an interactive that guides you through
the entire configuration process. Just run:

```
git town init
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

- GitHub: uses the [GitHub CLI](https://cli.github.com). If you prefer not to
  install `gh`, you can also configure an
  [access token](preferences/github-token.md) and use Git Town's built-in GitHub
  integration.
- GitLab: uses the [GitLab CLI](https://gitlab.com/gitlab-org/cli/-/tree/main).
  Without `glab`, you can configure an
  [access token](preferences/gitlab-token.md) for Git Town's built-in GitLab
  support.
- Bitbucket: requires a [username](preferences/bitbucket-username.md) and
  [app password](preferences/bitbucket-app-password.md)
- Gitea: requires an [access token](preferences/gitea-token.md)
- Forgejo: requires an [access token](preferences/forgejo-token.md)
