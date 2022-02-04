Feature: the branch to kill has a deleted tracking branch

  Background:
    Given my repo has the feature branches "old" and "other"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE              |
      | old    | local, remote | old commit           |
      | other  | local, remote | other feature commit |
    And the "old" branch gets deleted on the remote
    And I am on the "old" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                    |
      | old    | git fetch --prune --tags   |
      |        | git add -A                 |
      |        | git commit -m "WIP on old" |
      |        | git checkout main          |
      | main   | git branch -D old          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE              |
      | other  | local, remote | other feature commit |
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
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
    And I am now on the "old" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE              |
      | old    | local         | old commit           |
      | other  | local, remote | other feature commit |
    And my workspace has the uncommitted file again
    And my repo now has its initial branches and branch hierarchy
