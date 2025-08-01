@messyoutput
Feature: Configure a different development remote

  Background:
    Given a Git repo with origin
    And an additional "fork" remote with URL "https://github.com/forked/repo"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS       |
      | welcome                     | enter      |
      | aliases                     | enter      |
      | main branch                 | enter      |
      | perennial branches          |            |
      | perennial regex             | enter      |
      | feature regex               | enter      |
      | contribution regex          | enter      |
      | observed regex              | enter      |
      | new branch type             | enter      |
      | unknown branch type         | enter      |
      | dev-remote                  | up enter   |
      | origin hostname             | enter      |
      | forge type                  | enter      |
      | sync feature strategy       | enter      |
      | sync perennial strategy     | enter      |
      | sync prototype strategy     | enter      |
      | sync upstream               | enter      |
      | sync tags                   | enter      |
      | share new branches          | enter      |
      | push hook                   | enter      |
      | ship strategy               | enter      |
      | ship delete tracking branch | enter      |
      | config storage              | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --unset git-town.main-branch         |
      | git config git-town.unknown-branch-type feature |
    And the configuration file is now:
      """
      # More info around this file at https://www.git-town.com/configuration-file

      [branches]
      main = "main"

      [create]
      new-branch-type = "feature"
      share-new-branches = "no"

      [hosting]
      dev-remote = "fork"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
      feature-strategy = "merge"
      perennial-strategy = "ff-only"
      prototype-strategy = "merge"
      push-hook = true
      tags = true
      upstream = true
      """
