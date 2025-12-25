# Preferences

You can see all preferences via the [config](commands/config.md) command and
change them via the [setup assistant](commands/init.md) or manually.
Configuration data exists on multiple levels:

1. Team-wide configuration settings go into the
   [configuration file](configuration-file.md). These settings apply to all Git
   Town users working on the respective repository.

2. Each developer can configure their preferred Git Town settings for all
   repositories on their machine using global Git metadata. These settings
   override (1). For example, if I always want to use the `rebase`
   [sync-feature-strategy](https://www.git-town.com/preferences/sync-feature-strategy.html)
   in all my repositories, I would run:

   ```wrap
   git config --global git-town.sync-feature-strategy rebase
   ```

3. User and repo specific configuration settings go into local Git metadata,
   which takes precedence over (1) and (2). For example, if I want `rebase` as
   the default strategy for all my repositories, except in the `foo` repo I want
   to use `merge`, I'd first configure the global setting in (2), and then run
   in the `foo` repo:

   ```wrap
   git config git-town.sync-feature-strategy merge
   ```

4. All config settings can also be overridden via environment variables. For
   example, to load your [GitHub token](preferences/github-token.md) from the
   1Password CLI:

   ```bash
   GIT_TOWN_GITHUB_TOKEN=$(op read op://development/GitHub/credentials/personal_token) git town config
   ```
