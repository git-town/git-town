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
      | default branch type         | enter |
      | feature regex               | enter |
      | hosting platform            | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | sync-tags                   | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | new-branch-type             | enter |
      | ship-strategy               | enter |
      | ship-delete-tracking-branch | enter |
      | save config to config file  | enter |

  Scenario: result
    Then Git Town runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git Town setting "new-branch-type" still doesn't exist
    And local Git Town setting "main-branch" still doesn't exist
    And local Git Town setting "perennial-branches" still doesn't exist
    And local Git Town setting "default-branch-type" still doesn't exist
    And local Git Town setting "feature-regex" still doesn't exist
    And local Git Town setting "hosting-platform" still doesn't exist
    And local Git Town setting "push-new-branches" still doesn't exist
    And local Git Town setting "push-hook" still doesn't exist
    And local Git Town setting "sync-feature-strategy" still doesn't exist
    And local Git Town setting "sync-perennial-strategy" still doesn't exist
    And local Git Town setting "sync-upstream" still doesn't exist
    And local Git Town setting "sync-tags" still doesn't exist
    And local Git Town setting "ship-strategy" still doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist
    And the configuration file is now:
      """
      # Git Town configuration file
      #
      # Run "git town config setup" to add additional entries
      # to this file after updating Git Town.

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

      # All branches whose name matches this regular expression
      # are also considered perennial branches.
      #
      # If you are not sure, leave this empty.
      perennial-regex = ""

      [create]

      # The "new-branch-type" setting determines which branch type Git Town
      # creates when you run "git town hack", "append", or "prepend".
      #
      # More info at https://www.git-town.com/preferences/new-branch-type.
      new-branch-type = "feature"

      # Should Git Town push the new branches it creates
      # immediately to origin even if they are empty?
      #
      # When enabled, you can run "git push" right away
      # but creating new branches is slower and
      # it triggers an unnecessary CI run on the empty branch.
      #
      # When disabled, many Git Town commands execute faster
      # and Git Town will create the missing tracking branch
      # on the first run of "git town sync".
      push-new-branches = false

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

      [ship]

      # Should "git town ship" delete the tracking branch?
      # You want to disable this if your code hosting platform
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      delete-tracking-branch = true

      # Which method should Git Town use to ship feature branches?
      #
      # Options:
      #
      # - api: merge the proposal on your code hosting platform via the code hosting API
      # - fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch
      # - squash-merge: in your local repo, squash-merge the feature branch into its parent branch
      #
      # All options update proposals of child branches and remove the shipped branch locally and remotely.
      strategy = "api"

      [sync]

      # How should Git Town synchronize feature branches?
      # Feature branches are short-lived branches cut from
      # the main branch and shipped back into the main branch.
      # Typically you develop features and bug fixes on them,
      # hence their name.
      feature-strategy = "merge"

      # How should Git Town synchronize perennial branches?
      # Perennial branches have no parent branch.
      # The only updates they receive are additional commits
      # made to their tracking branch somewhere else.
      perennial-strategy = "rebase"

      # The "push-hook" setting determines whether Git Town
      # permits or prevents Git hooks while pushing branches.
      # Hooks are enabled by default. If your Git hooks are slow,
      # you can disable them to speed up branch syncing.
      #
      # When disabled, Git Town pushes using the "--no-verify" switch.
      # More info at https://www.git-town.com/preferences/push-hook.
      push-hook = true

      # Should "git town sync" sync tags with origin?
      tags = true

      # Should "git town sync" also fetch updates from the upstream remote?
      #
      # If an "upstream" remote exists, and this setting is enabled,
      # "git town sync" will also update the local main branch
      # with commits from the main branch at the upstream remote.
      #
      # This is useful if the repository you work on is a fork,
      # and you want to keep it in sync with the repo it was forked from.
      upstream = true
      """

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" still doesn't exist
    And global Git setting "alias.diff-parent" still doesn't exist
    And global Git setting "alias.hack" still doesn't exist
    And global Git setting "alias.delete" still doesn't exist
    And global Git setting "alias.prepend" still doesn't exist
    And global Git setting "alias.propose" still doesn't exist
    And global Git setting "alias.rename" still doesn't exist
    And global Git setting "alias.repo" still doesn't exist
    And global Git setting "alias.set-parent" still doesn't exist
    And global Git setting "alias.ship" still doesn't exist
    And global Git setting "alias.sync" still doesn't exist
    And local Git Town setting "new-branch-type" still doesn't exist
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
    And local Git Town setting "ship-strategy" still doesn't exist
    And local Git Town setting "ship-delete-tracking-branch" still doesn't exist
