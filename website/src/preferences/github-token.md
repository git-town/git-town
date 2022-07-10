# github-token

```
git-town.github-token=<token>
```

To interact with the GitHub API in your name, Git Town needs a
[personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
with the `repo` scope. You can create one at
https://github.com/settings/tokens/new. When you have created the token, run
`git config git-town.github-token <token>` (where `<token>` gets replaced with
the content of your GitHub access token) inside your code repository to store it
in the Git Town configuration for the current repository.
