Feature: the branch to kill has a deleted tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | old   | feature | main   | local, origin |
      | other | feature | main   | local, origin |
    And the current branch is "old"
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
      | BRANCH | COMMAND                                          |
      | old    | git fetch --prune --tags                         |
      |        | git add -A                                       |
      |        | git commit -m "Committing WIP for git town undo" |
      |        | git checkout other                               |
      | other  | git branch -D old                                |
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
      | BRANCH | COMMAND                                                     |
      | other  | git branch old {{ sha 'Committing WIP for git town undo' }} |
      |        | git checkout old                                            |
      | old    | git reset --soft HEAD~1                                     |
    And the current branch is now "old"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | old    | local         | old commit   |
      | other  | local, origin | other commit |
    And the uncommitted file still exists
    And the initial branches and lineage exist
