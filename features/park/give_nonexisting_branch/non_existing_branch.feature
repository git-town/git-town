Feature: cannot park non-existing branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town park feature non-existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      there is no branch "non-existing"
      """
    And branch "feature" still has type "feature"
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
