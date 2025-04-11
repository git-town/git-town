# Working with Forked Repositories

Git Town fully supports fork-based workflows. After cloning your fork locally,
add an `upstream` remote pointing to the original repo you forked from. You can
do this using the [Git remote](https://git-scm.com/docs/git-remote) command:

```
git remote add upstream <Git URL>
```

Once set up, [git town sync](../commands/sync.md) will automatically pull in
updates from the `upstream` repository. When you're ready to submit your
changes, [git town propose](../commands/propose.md) creates pull requests from
your fork to the original project.
