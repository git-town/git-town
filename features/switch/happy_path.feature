@messyoutput
Feature: switch branches

  Scenario: switching to another branch
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
    And the current branch is "alpha"
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the current branch is now "beta"
