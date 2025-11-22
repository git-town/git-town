# Git Town configuration file

Git Town can be configured through a configuration file named **git-town.toml**,
**.git-town.toml**, or **.git-branches.toml**. To create one, run:

```
git town init
```

Here is an example configuration file with the default settings:

```toml
[branches]
main = "" # must be set by the user
contribution-regex = ""
default-type = "feature"
feature-regex = ""
observed-regex = ""
perennial-regex = ""
perennials = []

[create]
branch-prefix = ""
new-branch-type = "feature"
share-new-branches = "no"

[hosting]
dev-remote = "origin"
origin-hostname = "" # use the hostname in the origin URL
forge-type = "" # auto-detect

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
auto-sync = true
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "rebase"
push-hook = true
tags = true
upstream = true
```
