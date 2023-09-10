Feature: offline mode

  Background:
    Given setting "sync-strategy" is "rebase"
    And offline mode is enabled
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                   |
      | feature | git checkout main         |
      | main    | git rebase origin/main    |
      |         | git checkout feature      |
      | feature | git rebase origin/feature |
      |         | git rebase main           |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local main commit     |
      |         |          | local feature commit  |
      |         | origin   | origin feature commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git reset --hard {{ sha 'local feature commit' }} |
      |         | git checkout main                                 |
      | main    | git checkout feature                              |
    And the current branch is still "feature"
    And now the initial commits exist
    And the initial branches and hierarchy exist
