# git town config sync-feature-strategy <merge|rebase>

The _sync-feature-strategy_ configuration command displays or sets the strategy
to use when syncing feature branches.

### Arguments

- without an argument, displays the current sync-feature strategy
- with `merge`, set the sync-feature strategy to merge changes into your feature
  branches
- with `rebase`, set the sync-feature strategy to rebase your feature branches
  against their parents and remote counterparts
