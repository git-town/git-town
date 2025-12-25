# Offline mode

If you have no internet connection, certain Git Town commands that perform
network requests will fail. Enabling offline mode omits all network operations
and thereby keeps Git Town working.

This setting applies to all repositories on your local machine.

## set via CLI

To put Git Town into offline mode, run
[git town offline](../commands/offline.md).

## Git metadata

```wrap
git config --global git-town.offline <true|false>
```

## environment variable

You can configure offline mode by setting the `GIT_TOWN_OFFLINE` environment
variable.
