Feature: parking a feature branch

  Background:
    Given the current branch is a feature branch "branch"
    And an uncommitted file
    When I run "git-town park"

  @this
  Scenario: result
    Then it runs no commands
    And the current branch is still "branch"
    And branch "branch" is now parked

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "branch"
    And branch "branch" is now a feature branch
