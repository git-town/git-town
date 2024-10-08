Feature: cannot ship a feature branch that is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        |
      | feature | local, origin | feature commit | conflicting_file |
    And the current branch is "other"
    And branch "feature" is active in another worktree
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" is active in another worktree
      """
    And the current branch is still "other"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "other"
