Feature: swapping a feature branch with a branch that has merge commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-1 | local, origin | commit 1    |
      | branch-2 | local, origin | commit 2    |
      | branch-3 | local, origin | commit 3    |
    And I ran "git merge main --message merging" on branch "branch-1"
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-2 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap because branch branch-1 contains merge commits - please compress and try again
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
