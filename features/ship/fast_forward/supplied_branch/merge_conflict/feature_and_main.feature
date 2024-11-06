Feature: does not ship an unsynced feature branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                 |
      | main    | local    | main commit             |
      | feature | local    | unsynced feature commit |
    And the current branch is "other"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    And I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash                   |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git merge --abort           |
      |        | git checkout other          |
      | other  | git stash pop               |
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And no merge is in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is now "other"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist now
    And the initial branches and lineage exist now
