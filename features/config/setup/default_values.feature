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
      | dev-remote                  | enter |
      | hosting platform            | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-prototype-strategy     | enter |
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
    And local Git Town setting "dev-remote" still doesn't exist
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
      # More info around this file at https://www.git-town.com/configuration-file

      [branches]
      main = "main"
      perennials = []
      perennial-regex = ""

      [create]
      new-branch-type = "feature"
      push-new-branches = false

      [hosting]
      dev-remote = "origin"
      # platform = ""
      # origin-hostname = ""

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      prototype-strategy = "merge"
      push-hook = true
      tags = true
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
    And local Git Town setting "dev-remote" still doesn't exist
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
