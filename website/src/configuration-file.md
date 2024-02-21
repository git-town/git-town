# Git Town Configuration File

Git Town can be configured through a configuration file with named
**.git-branches.toml**. To create one, execute:

```
git town config setup
```

Here is an example configuration file with the default settings:

```toml
push-new-branches = false
ship-delete-tracking-branch = true
sync-upstream = true

[branches]
main = ""             # must be set by the user
perennials = []
perennial-regex = ""

[hosting]
platform = ""         # auto-detect
origin-hostname = ""  # use the hostname in the origin URL

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
```
