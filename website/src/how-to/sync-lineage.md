# Sync branch lineage across team members or repository clones

Teams using Git Town, as well as individual users working across multiple
machines or repository clones, often need a way to share branch lineage.

Git Town supports this via proposals. When you [sync](../commands/sync.md) a
branch and Git Town doesn't know its parent, it checks your forge for an
existing proposal. If it finds one, Git Town uses the proposal's target branch
as the parent.

To automate proposal creation for new branches, set
[`share-new-branches`](https://www.git-town.com/preferences/share-new-branches.html)
to `propose`.
