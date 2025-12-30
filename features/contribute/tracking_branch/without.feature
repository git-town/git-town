Feature: make another local feature branch a contribution branch

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | local | feature | main   | local     |
    When I run "git-town contribute local"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch local is local only - branches you want to contribute to must have a remote branch because they are per definition other people's branches
      """
    And branch "local" still has type "feature"
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
