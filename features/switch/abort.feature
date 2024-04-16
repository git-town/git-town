Feature: switch branches

  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    When I run "git-town switch" and enter into the dialogs:
      | KEYS     |
      | down esc |
    Then it runs no commands
    And the current branch is still "alpha"
