Feature: making a parked branch a feature branch

  Background:
    Given the current branch is a parked branch "parked"
    When I run "git-town hack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | parked | git fetch --prune --tags |
    And it prints:
      """
      branch "parked" is now a feature branch
      """
    And branch "parked" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "parked" is now parked
