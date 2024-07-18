Feature: cannot make the current perennial branch a feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    When I run "git-town hack"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "perennial" is a perennial branch
      """
    And branch "perennial" is still perennial

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "perennial" is still perennial
