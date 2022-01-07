# Ship command

```
git ship [branch name] [-m|--message <message>]
```

The ship command ("let's ship these features") merge a completed feature branch
into the main branch and then removes the feature branch. Before the merge it
[syncs](sync.md) the branch to be shipped. After the merge it pushes the main
branch.

Opens the default Git editor with a pre-populated commit message that the user
can modify unless the `--message` parameter specifies it. To abort the ship,
submit an empty commit message.

This command ships only direct children of the main branch. To ship a nested
feature branch, you need to ship or kill all its ancestor branches first.

If you use GitHub or Gitea, have enabled
[API access to your hosting provider](../configure.md#enable-api-access-to-your-hosting-provider),
and the branch to be shipped has an open pull request, this command merges pull
requests via the API of the hosting service.

If your origin server deletes shipped branches, for example
[GitHub's feature to automatically delete head branches](https://help.github.com/en/github/administering-a-repository/managing-the-automatic-deletion-of-branches),
you can
[disable deleting remote branches](../configure.md#disable-deleting-remote-branches).
