Feature: park the current perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS |
      | perennial | perennial | local     |
    And the current branch is "perennial"
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot park perennial branches
      """
    And the perennial branches are still "perennial"
    And branch "perennial" still has type "perennial"
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
