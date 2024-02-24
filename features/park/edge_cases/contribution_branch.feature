Feature: parking a contribution branch

  Background:
    Given the current branch is a contribution branch "contribution"
    And an uncommitted file
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And the current branch is still "contribution"
    And branch "contribution" is still a contribution branch
    And the uncommitted file still exists
    And it prints the error:
      """
      cannot park contribution branches
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "contribution"
    And the uncommitted file still exists
    And branch "contribution" is still a contribution branch
