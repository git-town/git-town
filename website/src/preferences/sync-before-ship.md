# sync-before-ship

```
git-town.sync-before-ship=<true|false>
```

If set to `true`, [git ship](../commands/ship.md) executes
[git sync](../commands/sync.md) before shipping a branch.

Syncing before shipping allows you to deal with merge conflicts and the
resulting breakage on the feature branch instead of the main branch. This helps
keep the main branch green. The downside of syncing before shipping is that this
will trigger another CI run. This might block shipping until the CI run is
finished.
