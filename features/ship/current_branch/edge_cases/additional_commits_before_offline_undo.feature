Feature: undoing an offline ship with additional commits to main

  Background:
    Given offline mode is enabled
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    When I run "git-town ship -m 'feature done'"
    And I add commit "additional commit" to the "main" branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      cannot reset branch "main"
      """
    And it prints the error:
      """
      it received additional commits in the meantime
      """
    And the current branch is now "main"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE           |
      | main    | local    | feature done      |
      |         |          | additional commit |
      | feature | origin   | feature commit    |
    And the branches are now
      | REPOSITORY | BRANCHES      |
      | local      | main          |
      | origin     | main, feature |
    And the initial lineage exists
