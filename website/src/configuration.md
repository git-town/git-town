# Configuration

If your repository already contains a `.git-branches.toml` file, you are good to
go. If not or something doesn't work, run Git Town's setup assistant. It walks
you through every configuration option and gives you a chance to adjust it.

```
git town config setup
```

More information about the configuration file including how to create one
manually is [here](configuration-file.md).

### Access Tokens

API access multiplies Git Town's utility:

- if the parent of a branch is not known, Git Town can look for a pull requests
  of this branch and uses their parent branch
- updates affected pull requests when you prepend, rename, remove branches or
  change their parent
- click the "merge" button on a pull request from your CLI

Configuring API access takes only one minute. Here is how you do it:

- GitHub: [access token](preferences/github-token.md)
- GitLab: [access token](preferences/gitlab-token.md)
- Bitbucket: [username](preferences/bitbucket-username.md) and
  [app password](preferences/bitbucket-app-password.md)
- gitea: [access token](preferences/gitea-token.md)

### Install shell autocompletion

To have your shell auto-complete Git Town commands, set up
[shell autocompletion](commands/completions.md)
