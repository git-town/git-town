Feature: delete the current observed branch

  Background:
    Given the current branch is a observed branch "observed"
    And a feature branch "feature"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | feature  | local, origin | feature commit  |
      | observed | local, origin | observed commit |
    And an uncommitted file
    And the current branch is "observed" and the previous branch is "feature"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                         |
      | observed | git fetch --prune --tags        |
      |          | git add -A                      |
      |          | git commit -m "WIP on observed" |
      |          | git checkout feature            |
      | feature  | git branch -D observed          |
    And the current branch is now "feature"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY | BRANCHES                |
      | local      | main, feature           |
      | origin     | main, feature, observed |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | feature  | local, origin | feature commit  |
      | observed | origin        | observed commit |
    And this branch lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                         |
      | feature  | git branch observed {{ sha 'WIP on observed' }} |
      |          | git checkout observed                           |
      | observed | git reset --soft HEAD~1                         |
    And the current branch is now "observed"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
    And branch "observed" is now observed
