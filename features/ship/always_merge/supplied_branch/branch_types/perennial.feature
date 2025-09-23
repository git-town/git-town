Feature: does not ship perennial branches using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT | LOCATIONS     |
      | production | perennial |        | local, origin |
    And the current branch is "main"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship production"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot ship perennial branches
      """
    And no lineage exists now
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
