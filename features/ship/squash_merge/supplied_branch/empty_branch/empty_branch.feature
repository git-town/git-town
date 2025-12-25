Feature: does not ship empty feature branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | empty | feature | main   | local     |
      | other | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main   | local    | main commit    | common_file | common content |
      | empty  | local    | feature commit | common_file | common content |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "other"
    When I run "git-town ship empty"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
