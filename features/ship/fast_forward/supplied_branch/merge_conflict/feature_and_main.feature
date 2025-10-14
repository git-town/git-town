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
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "other"
    And I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | other  | git fetch --prune --tags    |
      |        | git checkout main           |
      | main   | git merge --ff-only feature |
      |        | git merge --abort           |
      |        | git checkout other          |
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And no merge is now in progress
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
