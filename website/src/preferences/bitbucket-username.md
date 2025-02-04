# Bitbucket username

Git Town can interact with Bitbucket Cloud and Bitbucket Data Center in your
name, for example to update pull requests as branches get created, shipped, or
deleted. To do so, Git Town needs your Bitbucket username and a
[Bitbucket app password](bitbucket-app-password.md).

## config file

Since usernames are user specific, you cannot add them to the config file.

## Git metadata

You can configure the Bitbucket username manually by running:

```wrap
git config [--global] git-town.bitbucket-username <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
