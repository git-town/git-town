Feature: making a contribution branch a feature branch

  Background:
    Given the current branch is a contribution branch "contribution"
    When I run "git-town hack"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And it prints:
      """
      branch "contribution" is now a feature branch
      """
    And branch "contribution" is now a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "contribution" is now a contribution branch
