Feature: parking the current contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME   | TYPE         | LOCATIONS |
      | branch | contribution | local     |
    And the current branch is "branch"
    When I run "git-town park"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "branch" is now parked
      """
    And the current branch is still "branch"
    And branch "branch" is now parked
    And there are now no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND       |
      | branch | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And the current branch is still "branch"
    And branch "branch" is now a contribution branch
    And there are now no parked branches
