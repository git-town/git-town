@messyoutput
Feature: override an existing Git alias

  Background:
    Given a Git repo with origin
    And I ran "git config --global alias.append checkout"
    And global Git setting "git-town.unknown-branch-type" is "feature"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG             | KEYS       |
      | welcome            | enter      |
      | aliases            | o enter    |
      | main branch        | enter      |
      | perennial branches |            |
      | perennial regex    | enter      |
      | origin hostname    | enter      |
      | forge type         | enter      |
      | enter all          | enter      |
      | config storage     | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                        |
      | git config --global alias.append "town append" |
      | git config --unset git-town.main-branch        |
    And global Git setting "alias.append" is now "town append"

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" is now "checkout"
    And local Git setting "git-town.main-branch" is now "main"
