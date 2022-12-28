# git town config sync-strategy <merge|rebase>

The _sync-strategy_ configuration command displays or sets your sync strategy.
The sync strategy specifies which strategy to use when syncing feature branches.

### Variations

- without an argument, displays the current sync strategy
- with `merge`, set the sync strategy to merge changes into your feature
  branches
- with `rebase`, set the sync strategy to rebase your feature branches against
  their parents and remote counterparts
