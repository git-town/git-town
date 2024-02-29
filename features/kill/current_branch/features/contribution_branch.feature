Feature: delete the current contribution branch

  Background:
    Given the current branch is a contribution branch "contribution"
    And a feature branch "feature"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
      | feature      | local, origin | feature commit      |
    And an uncommitted file
    And the current branch is "contribution" and the previous branch is "feature"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                             |
      | contribution | git fetch --prune --tags            |
      |              | git add -A                          |
      |              | git commit -m "WIP on contribution" |
      |              | git checkout feature                |
      | feature      | git branch -D contribution          |
    And the current branch is now "feature"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY | BRANCHES                    |
      | local      | main, feature               |
      | origin     | main, contribution, feature |
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | origin        | contribution commit |
      | feature      | local, origin | feature commit      |
    And this branch lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                                 |
      | feature      | git branch contribution {{ sha 'WIP on contribution' }} |
      |              | git checkout contribution                               |
      | contribution | git reset --soft HEAD~1                                 |
    And the current branch is now "contribution"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
    And branch "contribution" is now a contribution branch
