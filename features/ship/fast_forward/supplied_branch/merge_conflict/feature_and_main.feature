Feature: handle conflicts between the supplied feature branch and the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "other"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    And I run "git-town ship feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git add -A                  |
      |        | git stash                   |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git merge --abort           |
      |        | git checkout other          |
      | other  | git stash pop               |
    And it prints the error:
      """
      aborted because merge exited with error
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And no merge is in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is now "other"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial branches and lineage exist
