Feature: sync the current feature branch using the "rebase" sync strategy

  Background:
    Given setting "sync-strategy" is "rebase"
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
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git rebase origin/main      |
      |         | git push                    |
      |         | git checkout feature        |
      | feature | git rebase origin/feature   |
      |         | git rebase main             |
      |         | git push --force-with-lease |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local, origin | origin main commit    |
      |         |               | local main commit     |
      |         |               | origin feature commit |
      |         |               | local feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git checkout main    |
      | main    | git checkout feature |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local, origin | origin main commit    |
      |         |               | local main commit     |
      |         |               | origin feature commit |
      |         |               | local feature commit  |
    And the initial branches and hierarchy exist
