Feature: delete a parent branch

  Background:
    Given a feature branch "alpha"
    And a feature branch "beta" as a child of "alpha"
    And a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "gamma"
    And an uncommitted file
    When I run "git-town kill beta"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | gamma  | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git push origin :beta    |
      |        | git branch -D beta       |
      |        | git stash pop            |
    And the current branch is now "gamma"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES           |
      | local, origin | main, alpha, gamma |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | gamma  | local, origin | gamma commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | gamma  | alpha  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | gamma  | git add -A                              |
      |        | git stash                               |
      |        | git branch beta {{ sha 'beta commit' }} |
      |        | git push -u origin beta                 |
      |        | git stash pop                           |
    And the current branch is now "gamma"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
