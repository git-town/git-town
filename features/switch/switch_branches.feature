Feature: switch branches

  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    When I run "git-town switch" and enter into the dialogs:
      | DIALOG  | KEYS       |
      | welcome | down enter |
    Then it runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the current branch is now "beta"
