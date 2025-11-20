# Sync branch lineage across team members or repo clones

Teams that use Git Town, or users who use different computers or clones of the
same repo, need to share the branch lineage across different repository clones.

Git Town supports this through proposals. When you [sync](../commands/sync.md) a
branch, and Git Town doesn't know its parent branch, it checks whether your
forge has a proposal for this branch. If yes, it uses the target branch of the
proposal as the parent branch.

You can automate creation of proposals by setting
[share-new-branches](https://www.git-town.com/preferences/share-new-branches.html)
to `propose'.
