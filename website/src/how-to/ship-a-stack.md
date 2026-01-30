## Ship a branch from a stack without re-running CI

When you merge the proposal for a branch from a [stack](../stacked-changes.md),
you normally need to [sync](../commands/sync.md) afterward to propagate the
shipped changes through the rest of the stack before you can ship the next
branch in the stack. Doing so updates all branches in the stack, which triggers
another CI run for them. This is unnecessary because the actual code hasn't
changed.

You can avoid these unnecessary CI runs by shipping using Git Town's
[fast-forward strategy](../preferences/ship-strategy.md#fast-forward).

```bash
git switch <oldest-branch-in-stack>
git ship --strategy=fast-forward
git switch <next-branch-in-stack>
git ship --strategy=fast-forward
...
```

When fast-forward shipping, Git Town advances the main branch pointer to contain
the commits from the shipped branch. Since no commits are changed, downstream
branches remain valid and don't need to be re-synced and CI doesn't rerun.

This only works if the stack is in sync with the main branch, i.e. the main
branch hasn't received any new commits since you last synced the stack. If main
has new commits, a fast-forward is not possible and you must ship normally.
