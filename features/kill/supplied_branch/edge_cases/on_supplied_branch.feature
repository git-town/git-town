Feature: delete the current branch

  Background:
    Given my repo has the feature branches "other" and "current"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, remote | current commit |
      | other   | local, remote | other commit   |
    And I am on the "current" branch
    And my workspace has an uncommitted file
    When I run "git-town kill current"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                        |
      | current | git fetch --prune --tags       |
      |         | git push origin :current       |
      |         | git add -A                     |
      |         | git commit -m "WIP on current" |
      |         | git checkout main              |
      | main    | git branch -D current          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, remote | other commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch current {{ sha 'WIP on current' }} |
      |         | git checkout current                          |
      | current | git reset {{ sha 'current commit' }}          |
      |         | git push -u origin current                    |
    And I am now on the "current" branch
    And my workspace has the uncommitted file again
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
