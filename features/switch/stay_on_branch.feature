Feature: stay on the same branch

  @this
  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    When I run "git-town switch" and enter into the dialogs:
      | DIALOG | KEYS  |
      |        | enter |
    Then it runs the commands
      | BRANCH | COMMAND |
    And the current branch is still "alpha"
