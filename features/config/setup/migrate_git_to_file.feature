@messyoutput
Feature: migrate existing configuration in Git metadata to a config file

  Background:
    Given a Git repo with origin
    And the main branch is "main"
    And local Git setting "git-town.perennial-regex" is "release-.*"
    And local Git setting "git-town.perennial-branches" is "qa"
    And local Git setting "git-town.feature-regex" is "user-.*"
    And local Git setting "git-town.default-branch-type" is "observed"
    And local Git setting "git-town.dev-remote" is "fork"
    And local Git setting "git-town.push-new-branches" is "false"
    And local Git setting "git-town.push-hook" is "true"
    And local Git setting "git-town.new-branch-type" is "prototype"
    And local Git setting "git-town.ship-strategy" is "squash-merge"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    And local Git setting "git-town.sync-feature-strategy" is "merge"
    And local Git setting "git-town.sync-perennial-strategy" is "rebase"
    And local Git setting "git-town.sync-upstream" is "true"
    And local Git setting "git-town.sync-tags" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                               | KEYS  |
      | welcome                                   | enter |
      | add all aliases                           | enter |
      | accept the already configured main branch | enter |
      | perennial branches                        | enter |
      | perennial regex                           | enter |
      | default branch type                       | enter |
      | feature regex                             | enter |
      | dev-remote                                | enter |
      | hosting service                           | enter |
      | origin hostname                           | enter |
      | sync-feature-strategy                     | enter |
      | sync-perennial-strategy                   | enter |
      | sync-prototype-strategy                   | enter |
      | sync-upstream                             | enter |
      | sync-tags                                 | enter |
      | enable push-new-branches                  | enter |
      | disable the push hook                     | enter |
      | new-branch-type                           | enter |
      | ship-strategy                             | enter |
      | ship-delete-tracking-branch               | enter |
      | save config to config file                | enter |

  Scenario: result
    Then Git Town runs no commands
    And the main branch is now not set
    And there are now no perennial branches
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" now doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" now doesn't exist
    And local Git setting "git-town.sync-upstream" now doesn't exist
    And local Git setting "git-town.sync-tags" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.feature-regex" is still "user-.*"
    And local Git setting "git-town.default-branch-type" is still "observed"
    And local Git setting "git-town.push-new-branches" now doesn't exist
    And local Git setting "git-town.push-hook" now doesn't exist
    And local Git setting "git-town.new-branch-type" now doesn't exist
    And local Git setting "git-town.ship-strategy" now doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" now doesn't exist
    And the configuration file is now:
      """
      # More info around this file at https://www.git-town.com/configuration-file

      [branches]
      main = "main"
      perennials = ["qa"]
      perennial-regex = "release-.*"

      [create]
      new-branch-type = "prototype"
      push-new-branches = false

      [hosting]
      dev-remote = "origin"
      # platform = ""
      # origin-hostname = ""

      [ship]
      delete-tracking-branch = false
      strategy = "squash-merge"

      [sync]
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      prototype-strategy = "merge"
      push-hook = true
      tags = false
      upstream = true
      """

  Scenario: undo
    When I run "git-town undo"
    Then the main branch is now "main"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.new-branch-type" is now "prototype"
    And local Git setting "git-town.perennial-regex" is now "release-.*"
    And local Git setting "git-town.feature-regex" is now "user-.*"
    And local Git setting "git-town.default-branch-type" is now "observed"
    And local Git setting "git-town.push-new-branches" is now "false"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.ship-strategy" is now "squash-merge"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
    And local Git setting "git-town.sync-feature-strategy" is now "merge"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-upstream" is now "true"
    And local Git setting "git-town.sync-tags" is now "false"
