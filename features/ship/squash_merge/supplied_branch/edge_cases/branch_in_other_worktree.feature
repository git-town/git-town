Feature: shipping a branch that is checked out in another worktree

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
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship feature" and enter "feature done" for the commit message

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" is active in another worktree
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
