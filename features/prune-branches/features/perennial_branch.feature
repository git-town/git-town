Feature: remove perennial branch configuration when pruning a perennial branch

  Background:
    Given my repo has the perennial branches "active" and "old"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, remote | active commit |
      | old    | local, remote | old commit    |
    And the "old" branch gets deleted on the remote
    And I am on the "old" branch
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -D old        |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES     |
      | local, remote | main, active |
    And the perennial branches are now "active"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
    And I am now on the "old" branch
    And my repo now has the initial branches
    And the perennial branches are now "active" and "old"
