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
      | DIALOG        | KEYS     |
      | switch-branch | down esc |
    Then Git Town runs no commands
