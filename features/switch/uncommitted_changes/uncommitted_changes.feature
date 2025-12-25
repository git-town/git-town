@messyoutput
Feature: switch to another branch with uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | beta   | local, origin | beta commit |
    And the current branch is "alpha"
    And an uncommitted file
    When I run "git-town switch" and enter into the dialogs:
      | DIALOG        | KEYS       |
      | switch-branch | down enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the uncommitted file still exists
