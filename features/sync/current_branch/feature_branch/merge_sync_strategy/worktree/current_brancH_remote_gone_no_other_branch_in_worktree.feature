Feature: sync a branch whose branch is gone while main is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And origin deletes the "feature-1" branch
    And the current branch is "feature-1"
    And branch "main" is active in another worktree
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      no branch to switch to available
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
