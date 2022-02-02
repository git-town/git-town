Feature: cannot ship with uncommitted changes

  Background:
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town ship"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And my workspace still contains my uncommitted file

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And my workspace still contains my uncommitted file
