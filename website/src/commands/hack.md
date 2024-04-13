# git hack &lt;branch&gt;

The _hack_ command ("let's start hacking") creates a new feature branch with the
given name off the [main branch](../preferences/main-branch.md) and brings all
uncommitted changes over to it. When running without uncommitted changes in your
workspace, it also [syncs](sync.md) the main branch to ensure you develop on top
of the current state of the repository.

### Configuration

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

If [push-new-branches](../preferences/push-new-branches.md) is set, `git hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git hack` run fast. The first run of `git sync`
will create the remote tracking branch.
