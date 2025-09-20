@skipWindows
Feature: abort the ship by empty commit message

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "squash-merge"
    When I run "git-town ship" and enter an empty commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit                      |
      |         | git reset --hard                |
      |         | git checkout feature            |
    And Git Town prints the error:
      """
      aborted because merge exited with error
      """
    And the initial commits exist now
    And the initial branches and lineage exist now
#
# NOTE: Cannot test undo here.
# The Git Town command under test has not created an undoable runstate.
# Executing "git town undo" would undo the Git Town command executed during setup.
