# work with forked repositories

Git Town fully supports working with
[forked](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo)
repositories. After cloning your forked repo onto your local machine, add an
additional [Git remote](https://git-scm.com/docs/git-remote) called `upstream`
that points to the original repository (that you forked from).

Now [git town sync](../commands/sync.md) will pull in updates from the
`upstream` repository, and [git town propose](../commands/propose.md) will
create proposals from your fork to the original repository.
