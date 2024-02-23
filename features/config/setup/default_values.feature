Feature: Accepting all default values leads to a working setup

  Background:
    Given the branches "dev" and "production"
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main development branch     | enter |
      | perennial branches          | enter |
      | perennial regex             | enter |
      | hosting platform            | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | ship-delete-tracking-branch | enter |
      | sync-before-ship            | enter |
      | save config to config file  | enter |

  Scenario: result
    Then it runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git Town setting "hosting-platform" still doesn't exist
    And local Git Town setting "push-new-branches" is still not set
    And local Git Town setting "push-hook" is still not set
    And local Git Town setting "sync-feature-strategy" is still not set
    And local Git Town setting "sync-perennial-strategy" is still not set
    And local Git Town setting "sync-upstream" is still not set
    And local Git Town setting "ship-delete-tracking-branch" is still not set
    And local Git Town setting "sync-before-ship" is still not set
    And the configuration file is now:
      """
      # Git Town configuration file
      #
      # Run "git town config setup" to add additional entries
      # to this file after updating Git Town.
      #
      # The "push-hook" setting determines whether Git Town
      # permits or prevents Git hooks while pushing branches.
      # Hooks are enabled by default. If your Git hooks are slow,
      # you can disable them to speed up branch syncing.
      #
      # When disabled, Git Town pushes using the "--no-verify" switch.
      # More info at https://www.git-town.com/preferences/push-hook.
      push-hook = true

      # Should Git Town push the new branches it creates
      # immediately to origin even if they are empty?
      #
      # When enabled, you can run "git push" right away
      # but creating new branches is slower and
      # it triggers an unnecessary CI run on the empty branch.
      #
      # When disabled, many Git Town commands execute faster
      # and Git Town will create the missing tracking branch
      # on the first run of "git sync".
      push-new-branches = false

      # Should "git ship" delete the tracking branch?
      # You want to disable this if your code hosting platform
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      ship-delete-tracking-branch = true

      # Should "git ship" sync branches before shipping them?
      #
      # Guidance: enable when shipping branches locally on your machine
      # and disable when shipping feature branches via the code hosting
      # API or web UI.
      #
      # When enabled, branches are always fully up to date when shipped
      # and you get a chance to resolve merge conflicts
      # between the feature branch to ship and the main development branch
      # on the feature branch. This helps keep the main branch green.
      # But this also triggers another CI run and delays shipping.
      sync-before-ship = false

      # Should "git sync" also fetch updates from the upstream remote?
      #
      # If an "upstream" remote exists, and this setting is enabled,
      # "git sync" will also update the local main branch
      # with commits from the main branch at the upstream remote.
      #
      # This is useful if the repository you work on is a fork,
      # and you want to keep it in sync with the repo it was forked from.
      sync-upstream = true

      [branches]

      # The main branch is the branch from which you cut new feature branches,
      # and into which you ship feature branches when they are done.
      # This branch is often called "main", "master", or "development".
      main = "main"

      # Perennial branches are long-lived branches.
      # They are never shipped and have no ancestors.
      # Typically, perennial branches have names like
      # "development", "staging", "qa", "production", etc.
      #
      # See also the "perennial-regex" setting.
      perennials = []

      # All branches whose names match this regular expression
      # are also considered perennial branches.
      #
      # If you are not sure, leave this empty.
      perennial-regex = ""

      [hosting]

      # Knowing the type of code hosting platform allows Git Town
      # to open browser URLs and talk to the code hosting API.
      # Most people can leave this on "auto-detect".
      # Only change this if your code hosting server uses as custom URL.
      # platform = ""

      # When using SSH identities, define the hostname
      # of your source code repository. Only change this
      # if the auto-detection does not work for you.
      # origin-hostname = ""

      [sync-strategy]

      # How should Git Town synchronize feature branches?
      # Feature branches are short-lived branches cut from
      # the main branch and shipped back into the main branch.
      # Typically you develop features and bug fixes on them,
      # hence their name.
      feature-branches = "merge"

      # How should Git Town synchronize perennial branches?
      # Perennial branches have no parent branch.
      # The only updates they receive are additional commits
      # made to their tracking branch somewhere else.
      perennial-branches = "rebase"
      """

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" now doesn't exist
    And global Git setting "alias.diff-parent" now doesn't exist
    And global Git setting "alias.hack" now doesn't exist
    And global Git setting "alias.kill" now doesn't exist
    And global Git setting "alias.prepend" now doesn't exist
    And global Git setting "alias.propose" now doesn't exist
    And global Git setting "alias.rename-branch" now doesn't exist
    And global Git setting "alias.repo" now doesn't exist
    And global Git setting "alias.set-parent" now doesn't exist
    And global Git setting "alias.ship" now doesn't exist
    And global Git setting "alias.sync" now doesn't exist
    And the main branch is now ""
    And there are now no perennial branches
    And local Git Town setting "hosting-platform" now doesn't exist
    And local Git Town setting "github-token" now doesn't exist
    And local Git Town setting "hosting-origin-hostname" now doesn't exist
    And local Git Town setting "sync-feature-strategy" now doesn't exist
    And local Git Town setting "sync-perennial-strategy" now doesn't exist
    And local Git Town setting "sync-upstream" now doesn't exist
    And local Git Town setting "perennial-regex" now doesn't exist
    And local Git Town setting "push-new-branches" now doesn't exist
    And local Git Town setting "push-hook" now doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" now doesn't exist
    And local Git Town setting "sync-before-ship" now doesn't exist
