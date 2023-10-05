Feature: remove parent info of children of deleted perennial branches

  Background:
    Given the perennial branches "old" and "other"
    And a feature branch "old1" as a child of "old"
    And a feature branch "old2" as a child of "old"
    And a feature branch "other1" as a child of "other"
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
      | REPOSITORY    | BRANCHES                        |
      | local, origin | main, old1, old2, other, other1 |
    And the perennial branches are now "other"
    And this branch lineage exists now
      | BRANCH | PARENT |
      | other1 | other  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
    And the current branch is now "old"
    And the initial branches exist
    And the perennial branches are now "active" and "old"
