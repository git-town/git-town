Feature: swapping a branch with its contribution parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE         | PARENT | LOCATIONS     |
      | parent  | contribution | main   | local, origin |
      | current | feature      | parent | local, origin |
    And the current branch is "current"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch "parent" is a contribution branch
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
