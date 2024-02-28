Feature: cannot make a feature branch a feature branch

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town hack"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" is already a feature branch
      """
    And branch "feature" is still a feature branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "feature" is now a feature branch
