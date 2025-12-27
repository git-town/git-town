Feature: merging a branch whose parent was deleted remotely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | beta-file | beta content |
    And origin deletes the "alpha" branch
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch alpha is not in sync with its parent, please run "git town sync" and try again
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
