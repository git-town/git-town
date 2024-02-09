# Git Town Configuration File

Git Town can be configured through a configuration file. To create one, execute:

```
git town config setup
```

This creates a file **.git-branches.toml**. Here is one with the default
settings:

```toml
push-new-branches = false
ship-delete-tracking-branch = true
sync-upstream = true

[branches]
main = ""             # must be set by user
perennials = []

[hosting]
platform = ""         # defaults to auto-detect
origin-hostname = ""  # defaults to using the hostname in the origin URL

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
```
