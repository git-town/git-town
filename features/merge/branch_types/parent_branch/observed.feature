Feature: does not merge branches whose parent is observed

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE     | PARENT | LOCATIONS     |
      | parent  | observed |        | local, origin |
      | current | feature  | main   | local         |
    And the current branch is "current"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot merge branch current because its parent branch main has no parent
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
