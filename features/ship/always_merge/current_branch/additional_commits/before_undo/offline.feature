Feature: partially undo an offline ship using the always-merge strategy after additional commits to main

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And offline mode is enabled
    And the current branch is "feature"
    And I run "git-town ship" and close the editor

  Scenario: undo
    When I add commit "additional commit" to the "main" branch
    And I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And Git Town prints:
      """
      cannot reset branch main
      """
    And Git Town prints:
      """
      it received additional commits in the meantime
      """
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | feature commit         |
      |         |          | Merge branch 'feature' |
      |         |          | additional commit      |
      | feature | origin   | feature commit         |
