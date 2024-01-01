# Configuration commands

Git Town prompts for required configuration information during usage. Git Town
stores its configuration data inside
[Git configuration data](https://git-scm.com/docs/git-config). You can store
configuration values in the local or global Git configuration depending on
whether you want to share config settings between repositories or not. To see
your entire Git configuration, run `git config -l`. To see only the Git Town
configuration entries, run `git config --get-regexp git-town`. The following
commands read and write the configuration entries for you so that you don't have
to run Git configuration commands manually:

- [git town config](commands/config.md) - display or update your Git Town
  configuration
- [git town config main-branch](commands/config-main-branch.md) - display/set
  the main development branch for the current repo
- [git town config push-new-branches](commands/config-push-new-branches.md) -
  configure whether to push new empty branches to origin
- [git town config offline](commands/config-offline.md) - enable/disable offline
  mode
- [git town config perennial-branches](commands/config-perennial-branches.md) -
  display or update the perennial branches for the current repo
- [git town config sync-perennial-strategy](commands/config-sync-perennial-strategy.md) -
  display or set the strategy to update perennial branches
- [git town config sync-feature-strategy](commands/config-sync-feature-strategy.md) -
  display or set the strategy to sync via merges or rebases
