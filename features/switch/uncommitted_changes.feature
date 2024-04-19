Feature: switch branches

  @this
  Scenario: switching to another branch
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta"
    And an uncommitted file
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it prints:
      """
      You have uncommitted changes
      """
    And it runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout beta |
    And the current branch is now "beta"
