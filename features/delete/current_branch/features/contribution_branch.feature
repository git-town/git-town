Feature: delete the current contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
      | feature      | feature      | main   | local, origin |
    And the current branch is "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
      | feature      | local, origin | feature commit      |
    And an uncommitted file
    And the current branch is "contribution" and the previous branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                   |
      | contribution | git fetch --prune --tags                                  |
      |              | git add -A                                                |
      |              | git commit -m "Committing open changes on deleted branch" |
      |              | git checkout feature                                      |
      | feature      | git branch -D contribution                                |
    And the current branch is now "feature"
    And no uncommitted files exist now
    And the branches are now
      | REPOSITORY | BRANCHES                    |
      | local      | main, feature               |
      | origin     | main, contribution, feature |
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | origin        | contribution commit |
      | feature      | local, origin | feature commit      |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                       |
      | feature      | git branch contribution {{ sha 'Committing open changes on deleted branch' }} |
      |              | git checkout contribution                                                     |
      | contribution | git reset --soft HEAD~1                                                       |
    And the current branch is now "contribution"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
    And branch "contribution" now has type "contribution"
