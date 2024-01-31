# gitlab-token

```
git-town.gitlab-token=<token>
```

To interact with the GitLab API in your name, Git Town needs a
[personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)
with `api` scope. After you created your token, run the
[setup assistant](../commands/config-setup.md) and enter it there.

Alternatively, you can set the token manually by running
`git config git-town.gitlab-token <token>` inside your code repository.

GitLab supports different
[merge methods](https://docs.gitlab.com/ee/user/project/merge_requests/methods/)
that may need additional configuration. With GitLab's default settings, Git Town
will still create a merge request while shipping. Because shipping will squash
the commits, GitLab's "merge commit" and "merge commit with semi-linear history"
will produce the same result, creating two commits (change plus merge commit).
To get a linear history, the project needs to be configured to use the
"fast-forward merge" method.
