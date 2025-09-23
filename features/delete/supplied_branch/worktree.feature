Feature: delete a branch that is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | good | feature | main   | local, origin |
      | dead | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And branch "dead" is active in another worktree
    When I run "git-town delete dead"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "dead" is active in another worktree
      """
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
