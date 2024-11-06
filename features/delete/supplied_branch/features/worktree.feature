Feature: delete a branch that is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | good | feature | main   | local, origin |
      | dead | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And branch "dead" is active in another worktree
    And an uncommitted file
    When I run "git-town delete dead"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "dead" is active in another worktree
      """
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "good"
    And the uncommitted file still exists
