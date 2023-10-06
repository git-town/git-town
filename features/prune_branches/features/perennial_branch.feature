Feature: remove perennial branch configuration when pruning a perennial branch

  Background:
    Given the perennial branches "active" and "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | active | local, origin | active commit |
    And origin deletes the "old" branch
    And the current branch is "old"
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | old    | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git branch -d old        |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, active |
    And the perennial branches are now "active"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git branch old {{ sha 'Initial commit' }} |
      |        | git checkout old                          |
    And the current branch is now "old"
    And the initial branches exist
    And the perennial branches are now "active" and "old"
