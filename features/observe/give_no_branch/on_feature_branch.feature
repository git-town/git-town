Feature: observing the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town observe"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
      branch "feature" is now an observed branch
      """
    And the current branch is still "feature"
    And branch "feature" is now observed
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND       |
      | feature | git add -A    |
      |         | git stash     |
      |         | git stash pop |
    And the current branch is still "feature"
    And there are now no observed branches
    And the uncommitted file still exists
