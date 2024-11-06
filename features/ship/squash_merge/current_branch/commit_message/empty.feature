@skipWindows
Feature: abort the ship by empty commit message

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And Git Town setting "ship-strategy" is "squash-merge"
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
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
