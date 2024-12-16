# Git Town Configuration File

Git Town can be configured through a configuration file with named
**.git-branches.toml**. To create one, execute:

```
git town config setup
```

Here is an example configuration file with the default settings:

```toml
[branches]
main = ""                           # must be set by the user
contribution-regex = ""
default-type = "feature"
feature-regex = ""
observed-regex = ""
perennial-regex = ""
perennials = []

[create]
new-branch-type = "feature"
push-new-branches = false

[hosting]
dev-remote = "origin"
origin-hostname = ""  # use the hostname in the origin URL
platform = ""         # auto-detect

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
feature-strategy = "merge"
perennial-strategy = "rebase"
push-hook = true
tags = true
upstream = true

[sync-strategy]
prototype-branche = "merge"
```
