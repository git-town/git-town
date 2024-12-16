@messyoutput
Feature: Configure a different development remote

  Background:
    Given a Git repo with origin
    And an additional "fork" remote with URL "https://github.com/forked/repo"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS     |
      | welcome                     | enter    |
      | aliases                     | enter    |
      | main branch                 | enter    |
      | perennial branches          | enter    |
      | perennial regex             | enter    |
      | default branch type         | enter    |
      | feature regex               | enter    |
      | dev-remote                  | up enter |
      | hosting platform            | enter    |
      | origin hostname             | enter    |
      | sync-feature-strategy       | enter    |
      | sync-perennial-strategy     | enter    |
      | sync-prototype-strategy     | enter    |
      | sync-upstream               | enter    |
      | sync-tags                   | enter    |
      | push-new-branches           | enter    |
      | push-hook                   | enter    |
      | new-branch-type             | enter    |
      | ship-strategy               | enter    |
      | ship-delete-tracking-branch | enter    |
      | save config to config file  | enter    |

  Scenario: result
    Then Git Town runs no commands
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
      dev-remote = "fork"
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
