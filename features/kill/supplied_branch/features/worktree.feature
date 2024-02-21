Feature: delete a branch that is active in another worktree

  Background:
    Given the feature branches "good" and "dead"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And branch "dead" is active in another worktree
    And an uncommitted file
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
    And it prints the error:
      """
      branch "dead" is active in another worktree
      """
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "good"
    And the uncommitted file still exists
