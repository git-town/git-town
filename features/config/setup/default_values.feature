@messyoutput
Feature: Accepting all default values leads to a working setup

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | LOCATIONS     |
      | dev        | (none) | local, origin |
      | production | (none) | local, origin |
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main branch                 | enter |
      | perennial branches          | enter |
      | perennial regex             | enter |
      | hosting platform            | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | sync-tags                   | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | create-prototype-branches   | enter |
      | ship-delete-tracking-branch | enter |
      | save config to config file  | enter |

  Scenario: result
    Then it runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git Town setting "create-prototype-branches" still doesn't exist
    And local Git Town setting "main-branch" still doesn't exist
    And local Git Town setting "perennial-branches" still doesn't exist
    And local Git Town setting "hosting-platform" still doesn't exist
    And local Git Town setting "push-new-branches" still doesn't exist
    And local Git Town setting "push-hook" still doesn't exist
    And local Git Town setting "sync-feature-strategy" still doesn't exist
    And local Git Town setting "sync-perennial-strategy" still doesn't exist
    And local Git Town setting "sync-upstream" still doesn't exist
    And local Git Town setting "sync-tags" still doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist
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

      # The "create-prototype-branches" setting determines whether Git Town
      # always creates prototype branches.
      # Prototype branches sync only locally and don't create a tracking branch
      # until they are proposed.
      #
      # More info at https://www.git-town.com/preferences/create-prototype-branches.
      create-prototype-branches = false

      # Should "git ship" delete the tracking branch?
      # You want to disable this if your code hosting platform
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      ship-delete-tracking-branch = true

      # Should "git sync" sync tags with origin?
      sync-tags = true

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
    Then global Git setting "alias.append" still doesn't exist
    And global Git setting "alias.diff-parent" still doesn't exist
    And global Git setting "alias.hack" still doesn't exist
    And global Git setting "alias.kill" still doesn't exist
    And global Git setting "alias.prepend" still doesn't exist
    And global Git setting "alias.propose" still doesn't exist
    And global Git setting "alias.rename-branch" still doesn't exist
    And global Git setting "alias.repo" still doesn't exist
    And global Git setting "alias.set-parent" still doesn't exist
    And global Git setting "alias.ship" still doesn't exist
    And global Git setting "alias.sync" still doesn't exist
    And local Git Town setting "create-prototype-branches" still doesn't exist
    And local Git Town setting "main-branch" still doesn't exist
    And local Git Town setting "perennial-branches" still doesn't exist
    And local Git Town setting "hosting-platform" still doesn't exist
    And local Git Town setting "github-token" still doesn't exist
    And local Git Town setting "hosting-origin-hostname" still doesn't exist
    And local Git Town setting "sync-feature-strategy" still doesn't exist
    And local Git Town setting "sync-perennial-strategy" still doesn't exist
    And local Git Town setting "sync-upstream" still doesn't exist
    And local Git Town setting "sync-tags" still doesn't exist
    And local Git Town setting "perennial-regex" still doesn't exist
    And local Git Town setting "push-new-branches" still doesn't exist
    And local Git Town setting "push-hook" still doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist
