Feature: hoisting a branch whose parent is checked out in another workspace

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-1 | local, origin | commit 1a |
      | branch-1 | local, origin | commit 1b |
    And the branches
      | NAME     | TYPE   | PARENT   | LOCATIONS     |
      | branch-2 | parked | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-2 | local, origin | commit 2a |
      | branch-2 | local, origin | commit 2b |
    And branch "branch-1" is active in another worktree
    And the current branch is "branch-2"
    When I run "git-town hoist"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot hoist because branch "branch-1" is active in another worktree
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the initial commits exist now
    And the initial lineage exists now
