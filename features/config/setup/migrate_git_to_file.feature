@messyoutput
Feature: migrate existing configuration in Git metadata to a config file

  Background:
    Given a Git repo with origin
    And the branch
      | NAME | TYPE      | LOCATIONS     |
      | qa   | perennial | local, origin |
    And the main branch is "main"
    And local Git Town setting "perennial-regex" is "release-.*"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "true"
    And local Git Town setting "create-prototype-branches" is "true"
    And local Git Town setting "ship-delete-tracking-branch" is "false"
    And local Git Town setting "sync-feature-strategy" is "merge"
    And local Git Town setting "sync-perennial-strategy" is "rebase"
    And local Git Town setting "sync-upstream" is "true"
    And local Git Town setting "sync-tags" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                               | KEYS  |
      | welcome                                   | enter |
      | add all aliases                           | enter |
      | accept the already configured main branch | enter |
      | perennial branches                        | enter |
      | perennial regex                           | enter |
      | hosting service                           | enter |
      | origin hostname                           | enter |
      | sync-feature-strategy                     | enter |
      | sync-perennial-strategy                   | enter |
      | sync-upstream                             | enter |
      | sync-tags                                 | enter |
      | enable push-new-branches                  | enter |
      | disable the push hook                     | enter |
      | create-prototype-branches                 | enter |
      | disable ship-delete-tracking-branch       | enter |
      | save config to config file                | enter |

  Scenario: result
    Then it runs no commands
    And the main branch is now not set
    And there are now no perennial branches
    And local Git Town setting "hosting-platform" now doesn't exist
    And local Git Town setting "hosting-origin-hostname" now doesn't exist
    And local Git Town setting "sync-feature-strategy" now doesn't exist
    And local Git Town setting "sync-perennial-strategy" now doesn't exist
    And local Git Town setting "sync-upstream" now doesn't exist
    And local Git Town setting "sync-tags" now doesn't exist
    And local Git Town setting "perennial-regex" now doesn't exist
    And local Git Town setting "push-new-branches" now doesn't exist
    And local Git Town setting "push-hook" now doesn't exist
    And local Git Town setting "create-prototype-branches" now doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" now doesn't exist
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
      create-prototype-branches = true

      # Should "git ship" delete the tracking branch?
      # You want to disable this if your code hosting platform
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      ship-delete-tracking-branch = false

      # Should "git sync" sync tags with origin?
      sync-tags = false

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
      perennials = ["qa"]

      # All branches whose names match this regular expression
      # are also considered perennial branches.
      #
      # If you are not sure, leave this empty.
      perennial-regex = "release-.*"

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
    Then the main branch is now "main"
    And local Git Town setting "create-prototype-branches" is now "true"
    And local Git Town setting "perennial-regex" is now "release-.*"
    And local Git Town setting "push-new-branches" is now "false"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "false"
    And local Git Town setting "sync-feature-strategy" is now "merge"
    And local Git Town setting "sync-perennial-strategy" is now "rebase"
    And local Git Town setting "sync-upstream" is now "true"
    And local Git Town setting "sync-tags" is now "false"
