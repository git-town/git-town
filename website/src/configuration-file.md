# Git Town Configuration File

Git Town can be configured through a configuration file with named
**.git-branches.toml**. To create one, execute:

```
git town config setup
```

Here is an example configuration file with the default settings:

```toml
push-new-branches = false

[branches]
main = ""             # must be set by the user
perennials = []
perennial-regex = ""

[hosting]
platform = ""         # auto-detect
origin-hostname = ""  # use the hostname in the origin URL

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
push-hook = true
tags = true
upstream = true

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
```
