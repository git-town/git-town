# Creating and shipping hotfixes

Hotfixes differ from regular changes in that they’re based on a different
[perennial branch](../preferences/perennial-branches.md)—typically something
like `production` or `staging`—and get merged back into that branch.

For example, to create a hotfix from the `production` branch:

```
git checkout production
git append my-hotfix
```

Now, when you run [git town sync](../commands/sync.md), it'll sync your
`my-hotfix` branch with `production` instead of the main branch. When you're
ready to submit the fix, [git town propose](../commands/propose.md) will create
a pull request from your hotfix branch back into `production`.
