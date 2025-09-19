Feature: disable syncing via CLI

  Background:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [sync]
      auto-sync = false
      """
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                |
      | main     | local    | local main commit      |
      | main     | origin   | origin main commit     |
      | branch-1 | origin   | origin branch-1 commit |
    And the current branch is "branch-1"
    When I run "git-town append branch-2"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | branch-1 | git checkout -b branch-2 |
    And the initial commits exist now
    And this lineage exists now
      """
      main
        branch-1
          branch-2
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | branch-2 | git checkout branch-1  |
      | branch-1 | git branch -D branch-2 |
    And the initial commits exist now
    And the initial lineage exists now
