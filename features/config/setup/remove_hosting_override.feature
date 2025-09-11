@messyoutput
Feature: remove an existing forge type override

  Background:
    Given a Git repo with origin
    And local Git setting "git-town.forge-type" is "github"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                    | DESCRIPTION                                 |
      | welcome            | enter                   |                                             |
      | aliases            | enter                   |                                             |
      | main branch        | enter                   |                                             |
      | perennial branches |                         | no input here since the dialog doesn't show |
      | origin hostname    | enter                   |                                             |
      | forge type         | up up up up up up enter |                                             |
      | enter all          | enter                   |                                             |
      | config storage     | enter                   |                                             |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                |
      | git config --unset git-town.forge-type |
    And local Git setting "git-town.forge-type" now doesn't exist

  Scenario: undo
    When I run "git-town undo"
    And local Git setting "git-town.forge-type" is now "github"
