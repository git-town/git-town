Feature: does not ship out-of-sync branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                |
      | feature | local    | unsynced local commit  |
      |         | origin   | unsynced origin commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" is not in sync
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
