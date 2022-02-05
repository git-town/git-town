Feature: delete branches that were shipped or removed on another machine

  Background:
    Given my repo has the feature branches "active" and "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, remote | active commit |
      | old    | local, remote | old commit    |
    And the "old" branch gets deleted on the remote
    And I am on the "old" branch
    And my workspace has an uncommitted file
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -D old        |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES     |
      | local, remote | main, active |
    And Git Town now knows about this branch hierarchy
      | BRANCH | PARENT |
      | active | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
    And I am now on the "old" branch
    And my workspace still contains my uncommitted file
    And my repo now has its initial branches and branch hierarchy
