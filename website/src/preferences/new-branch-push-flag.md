# new-branch-push-flag

```
git-town.new-branch-push-flag=<true|false>
```

By default, Git Town creates new feature branches only in your local repository.
Git Town doesn't push them to the `origin` remote because that makes creating
branches slower and triggers an unnecessary CI run for a branch containing no
changes. Running [git sync](../commands/sync.md) or
[git new-pull-request](../commands/new-pull-request.md) will push the branch to
origin. If you prefer to push new branches when creating them, set this option
to `true` by running:

```
git config [--global] new-branch-push-flag <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies only to the Git repo you are
currently in.
