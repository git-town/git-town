Feature: make a perennial branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE      | PARENT | LOCATIONS |
      | existing | perennial | main   | local     |
    And the current branch is "main"
    When I run "git-town feature existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make perennial branches feature branches
      """
    And the initial branches and lineage exist now
    And branch "existing" still has type "perennial"
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
