Feature: undo offline sync after additional commits to the feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And offline mode is enabled
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: add commit and undo
    When I add commit "additional commit" to the "feature" branch
    And I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      cannot reset branch feature
      """
    And Git Town prints:
      """
      because it received additional commits in the meantime
      """
    And the initial branches and lineage exist now
