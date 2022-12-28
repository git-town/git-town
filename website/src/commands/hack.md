# git hack &lt;branch&gt;

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch-name.md) and brings
all uncommitted changes over to it. Before it does that, it [syncs](sync.md) the
main branch to ensure commits into the new branch are on top of the current
state of the repository.

### Variations

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

If [new-branch-push-flag](config-new-branch-push-flag.md) is set, `git hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git hack` run fast. The first run of `git sync`
will create the remote tracking branch.
