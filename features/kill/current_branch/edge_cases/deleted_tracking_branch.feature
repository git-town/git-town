Feature: the branch to kill has a deleted tracking branch

  Background:
    Given the current branch is a feature branch "old"
    And a feature branch "other"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | old    | local, origin | old commit   |
      | other  | local, origin | other commit |
    And origin deletes the "old" branch
    And an uncommitted file
    And the current branch is "old" and the previous branch is "other"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                    |
      | old    | git fetch --prune --tags   |
      |        | git add -A                 |
      |        | git commit -m "WIP on old" |
      |        | git checkout other         |
      | other  | git branch -D old          |
    And the current branch is now "other"
    And no uncommitted files exist
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | other  | git branch old {{ sha 'WIP on old' }} |
      |        | git checkout old                      |
      | old    | git reset --soft HEAD~1               |
    And the current branch is now "old"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | old    | local         | old commit   |
      | other  | local, origin | other commit |
    And the uncommitted file still exists
    And the initial branches and lineage exist
