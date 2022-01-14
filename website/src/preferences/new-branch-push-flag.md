# new-branch-push-flag

```
git-town.new-branch-push-flag=<true|false>
```

By default, [git hack](../commands/hack.md), [append](../commands/append.md),
and [prepend](../commands/prepend.md) create new feature branches in your local
repository. They don't push the new branch to `origin` because that would be
slow and unnecessarily trigger a CI run for the empty branch. Git Town will push
the new feature branch the first time you run [git sync](../commands/sync.md).

If you prefer that new branches get pushed when creating them, enable the
`new-branch-push-flag` configuration setting by running the
[git town new-branch-push-flag](../commands/new-branch-push-flag.md) command.
