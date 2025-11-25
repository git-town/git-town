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

<!-- keep-sorted start -->

- [git town completions](commands/completions.md) - set up shell autocomplete
- [git town config](commands/config.md) - display or update your Git Town
  configuration
- [git town config get-parent](commands/config-get-parent.md) - display the name
  of the parent branch
- [git town config remove](commands/config-remove.md) - remove the Git Town
  configuration
- [git town offline](commands/offline.md) - enable/disable offline mode
- [git town init](commands/init.md) - setup assistant for all config settings

<!-- keep-sorted end -->
