# push-new-branches

```
git-town.push-new-branches=<true|false>
```

By default, Git Town does not push new feature branches to the `origin` remote
since that would make creating branches slower and triggers an unnecessary CI
run for a branch containing no changes. Running [git sync](../commands/sync.md)
or [git propose](../commands/propose.md) will push the branch to origin later.
If you prefer to push new branches upon creation, set this option to `true` by
running:

```
git config [--global] push-new-branches <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the Git repo you are in.
