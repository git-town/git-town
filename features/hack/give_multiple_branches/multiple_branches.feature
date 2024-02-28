Feature: making multiple branches a feature branch

  Background:
    Given a contribution branch "contribution"
    And an observed branch "observed"
    When I run "git-town hack contribution observed"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "contribution" is now a contribution branch
