Feature: cannot ship a branch without proposal via API

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "api"
    And the current branch is "feature"
    When I run "git-town ship -m done"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                        |
      | feature | git fetch --prune --tags                       |
      |         | Finding proposal from feature into main ... ok |
    And Git Town prints the error:
      """
      cannot ship branch feature via API because it has no proposal
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
