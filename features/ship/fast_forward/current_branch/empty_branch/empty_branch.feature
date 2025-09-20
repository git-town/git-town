Feature: does not ship an empty branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | empty | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME | FILE CONTENT |
      | main   | local    | main commit  | same_file | same content |
      | empty  | local    | empty commit | same_file | same content |
    And the current branch is "empty"
    And Git setting "git-town.ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | empty  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the initial commits exist now
    And the initial branches and lineage exist now
#
# NOTE: Cannot test undo here.
# The Git Town command under test has not created an undoable runstate.
# Executing "git town undo" would undo the Git Town command executed during setup.
