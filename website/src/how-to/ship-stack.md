## Ship several branches in a stack

After you merge the proposal for the oldest branch in a
[stack](../stacked-changes.md), you normally need to [sync](../commands/sync.md)
to propagate the shipped changes through the rest of the stack before you can
ship the next branch in the stack. This triggers a CI run for each branch in the
stack, which is unnecessary and wasteful because the code hasn't changed.

You can avoid these unnecessary CI runs by shipping using Git Town's
[fast-forward strategy](../preferences/ship-strategy.md#fast-forward):

```bash
git switch <oldest-branch-in-stack>
git ship --strategy=fast-forward
git switch <next-branch-in-stack>
git ship --strategy=fast-forward
...
```

With fast-forward shipping, Git Town
[fast-forwards](https://git-scm.com/docs/git-merge#_fast_forward_merge) the main
branch to include the commits from the shipped branch. Since no commits are
rewritten (only branch pointers move), downstream branches in the stack remain
in sync and CI doesn't get triggered.

This only works if the stack is in sync with the main branch, i.e. the main
branch hasn't received new commits since you last synced the stack. If main has
new commits, a fast-forward is no longer possible and you must either sync the
stack again or ship using a different strategy.
