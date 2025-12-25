Feature: cannot ship a branch without connector

  Background:
    Given a Git repo with origin
    And the origin is "git@unknown.com:something/whatever.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "api"
    And the current branch is "feature"
    And a proposal for this branch does not exist
    When I run "git-town ship -m done"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please configure API access to your forge, more info at https://www.git-town.com/configuration#access-tokens
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
