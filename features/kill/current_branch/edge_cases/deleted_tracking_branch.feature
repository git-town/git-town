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
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                    |
      | old    | git fetch --prune --tags   |
      |        | git add -A                 |
      |        | git commit -m "WIP on old" |
      |        | git checkout main          |
      | main   | git branch -D old          |
    And the current branch is now "main"
    And my repo doesn't have any uncommitted files
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git branch old {{ sha 'WIP on old' }} |
      |        | git checkout old                      |
      | old    | git reset {{ sha 'old commit' }}      |
    And the current branch is now "old"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | old    | local         | old commit   |
      | other  | local, origin | other commit |
    And my workspace has the uncommitted file again
    And the initial branches and hierarchy exist
