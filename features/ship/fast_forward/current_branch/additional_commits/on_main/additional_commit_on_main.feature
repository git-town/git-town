Feature: cannot ship not-up-to-date feature branches using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | main    | local, origin | main commit    |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git merge --abort           |
      |         | git checkout feature        |
    And it prints the error:
      """
      aborted because merge exited with error
      """
    And the current branch is still "feature"
    And the initial branches and lineage exist
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
