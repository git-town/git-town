Feature: detaching an empty branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS     |
      | staging | perennial | main   | local, origin |
    And the current branch is "staging"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | staging | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot detach branches without parent
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
