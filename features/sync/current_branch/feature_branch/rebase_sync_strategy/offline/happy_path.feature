Feature: offline mode

  Background:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And offline mode is enabled
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git checkout main      |
      | main    | git rebase origin/main |
      |         | git checkout feature   |
      | feature | git rebase main        |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local main commit     |
      |         |          | local feature commit  |
      |         | origin   | origin feature commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                      |
      | feature | git reset --hard {{ sha-before-run 'local feature commit' }} |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
