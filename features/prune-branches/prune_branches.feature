Feature: delete branches that were shipped or removed on another machine

  Background:
    Given the feature branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
      | old    | local, origin | old commit    |
    And origin deletes the "old" branch
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
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And Git Town is now aware of this branch hierarchy
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
    And the initial branch setup and hierarchy exist
