Feature: observe the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the current branch is "contribution"
    And an uncommitted file
    When I run "git-town observe"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      branch "contribution" is now an observed branch
      """
    And branch "contribution" is now observed
    And there are now no contribution branches
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND       |
      | contribution | git add -A    |
      |              | git stash     |
      |              | git stash pop |
    And the current branch is still "contribution"
    And branch "contribution" is now a contribution branch
    And there are now no observed branches
    And the uncommitted file still exists
