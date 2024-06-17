Feature: observing multiple branches

  Background:
    Given the feature branches "feature-1", "feature-2", and "feature-3"
    And an uncommitted file
    When I run "git-town observe feature-1 feature-2 feature-3"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature-1" is now an observed branch
      """
    And branch "feature-1" is now observed
    And it prints:
      """
      branch "feature-2" is now an observed branch
      """
    And branch "feature-2" is now observed
    And it prints:
      """
      branch "feature-3" is now an observed branch
      """
    And branch "feature-3" is now observed
    And the current branch is still "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | main   | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And there are now no observed branches
    And the current branch is still "main"
    And the uncommitted file still exists
