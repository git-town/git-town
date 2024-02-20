Feature: undo offline sync after additional commits to the feature branch

  Background:
    Given offline mode is enabled
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"
    And I add commit "additional commit" to the "feature" branch

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      cannot reset branch "feature"
      """
    And it prints:
      """
      because it received additional commits in the meantime
      """
    And it runs no commands
    And the current branch is still "feature"
    And the initial branches and lineage exist
