Feature: undo offline sync after additional commits to the feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And offline mode is enabled
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And the current branch is "feature"
    When I run "git-town sync"
    And I add commit "additional commit" to the "feature" branch

  Scenario: undo
    When I run "git-town undo"
    Then Git Town prints:
      """
      cannot reset branch "feature"
      """
    And Git Town prints:
      """
      because it received additional commits in the meantime
      """
    And Git Town runs no commands
    And the initial branches and lineage exist now
