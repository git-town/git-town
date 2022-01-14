# main-branch-name

```
git-town.main-branch-name=<branch>
```

The main branch is the parent branch for new feature branches created with
[git hack](../commands/hack.md) and the branch into which Git Town
[ships](../commands/ship.md) finished feature branches.

Git Town remembers the name of the main branch in the `main-branch-name`
setting. When unknown, Git Town automatically asks for it. You can run
[git town main-branch](../commands/main-branch.md) to see or update the
configured main branch.
