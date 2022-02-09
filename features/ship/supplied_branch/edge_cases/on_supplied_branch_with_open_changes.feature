Feature: does not ship a branch that has open changes

  Background:
    Given a feature branch "feature"
    And my workspace has an uncommitted file
    And I am on the "feature" branch
    When I run "git-town ship feature"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
