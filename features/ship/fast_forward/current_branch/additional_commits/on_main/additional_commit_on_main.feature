Feature: cannot ship not-up-to-date feature branches using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | main    | local, origin | main commit    |
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "feature"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git merge --abort           |
      |         | git checkout feature        |
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
