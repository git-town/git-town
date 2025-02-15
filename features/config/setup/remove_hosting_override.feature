@messyoutput
Feature: remove an existing forge type override

  Background:
    Given a Git repo with origin
    And local Git setting "git-town.forge-type" is "github"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | down enter        |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | default branch type         | enter             |                                             |
      | feature regex               | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | forge type                  | up up up up enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | enter             |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                |
      | git config --unset git-town.forge-type |
    And local Git setting "git-town.forge-type" now doesn't exist

  Scenario: undo
    When I run "git-town undo"
    And local Git setting "git-town.forge-type" is now "github"
