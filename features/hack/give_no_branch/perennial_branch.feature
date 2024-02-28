Feature: cannot make a perennial branch a feature branch

  Background:
    Given the current branch is a perennial branch "perennial"
    When I run "git-town hack"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And it prints the error:
      """
      branch "perennial" is a perennial branch
      """
    And branch "perennial" is still perennial

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "perennial" is still perennial
