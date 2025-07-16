@messyoutput
Feature: override an existing Git alias

  Background:
    Given a Git repo with origin
    And I ran "git config --global alias.append checkout"
    And local Git setting "git-town.unknown-branch-type" is "feature"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS    | DESCRIPTION                     |
      | welcome                     | enter   |                                 |
      | aliases                     | o enter |                                 |
      | main branch                 | enter   |                                 |
      | perennial branches          |         | skipped because only one branch |
      | perennial regex             | enter   |                                 |
      | feature regex               | enter   |                                 |
      | contribution regex          | enter   |                                 |
      | observed regex              | enter   |                                 |
      | unknown branch type         | enter   |                                 |
      | dev remote                  |         | skipped because only one remote |
      | origin hostname             | enter   |                                 |
      | forge type                  | enter   |                                 |
      | sync-feature-strategy       | enter   |                                 |
      | sync-perennial-strategy     | enter   |                                 |
      | sync-prototype-strategy     | enter   |                                 |
      | sync-upstream               | enter   |                                 |
      | sync-tags                   | enter   |                                 |
      | share-new-branches          | enter   |                                 |
      | push-hook                   | enter   |                                 |
      | new-branch-type             | enter   |                                 |
      | ship-strategy               | enter   |                                 |
      | ship-delete-tracking-branch | enter   |                                 |
      | save config to config file  | enter   |                                 |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --global alias.append "town append"  |
      | git config --unset git-town.main-branch         |
      | git config git-town.unknown-branch-type feature |
    And global Git setting "alias.append" is now "town append"
    And the configuration file is now:
      """
      # More info around this file at https://www.git-town.com/configuration-file
      
      [branches]
      main = "main"
      
      [create]
      new-branch-type = "feature"
      share-new-branches = "no"
      
      [ship]
      delete-tracking-branch = true
      strategy = "api"
      
      [sync]
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      tags = true
      upstream = true
      """

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" is now "checkout"
    And local Git setting "git-town.main-branch" is now "main"
