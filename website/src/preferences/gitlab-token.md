# gitlab-token

```
git-town.gitlab-token=<token>
```

To interact with the GitLab API in your name, Git Town needs a
[personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).
After you created your token, run `git config git-town.gitlab-token <token>`
inside your code repository to store it in the Git Town configuration for the
current repository.
