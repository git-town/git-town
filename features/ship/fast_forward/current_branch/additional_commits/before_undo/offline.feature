Feature: partially undo an offline ship using the fast-forward strategy after additional commits to main

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And offline mode is enabled
    And Git setting "git-town.ship-strategy" is "fast-forward"
    When I run "git-town ship"
    And I add commit "additional commit" to the "main" branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And Git Town prints:
      """
      cannot reset branch "main"
      """
    And Git Town prints:
      """
      it received additional commits in the meantime
      """
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE           |
      | main    | local    | feature commit    |
      |         |          | additional commit |
      | feature | origin   | feature commit    |
    And the initial branches and lineage exist now
