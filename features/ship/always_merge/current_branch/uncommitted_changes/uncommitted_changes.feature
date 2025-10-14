Feature: does not ship uncommitted changes using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And the uncommitted file still exists
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
