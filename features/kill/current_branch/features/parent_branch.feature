Feature: delete a branch within a branch chain

  Background:
    Given a feature branch "alpha"
    And a feature branch "beta" as a child of "alpha"
    And a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "beta"
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | beta   | git fetch --prune --tags    |
      |        | git push origin :beta       |
      |        | git add -A                  |
      |        | git commit -m "WIP on beta" |
      |        | git checkout alpha          |
      | alpha  | git branch -D beta          |
    And the current branch is now "alpha"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY    | BRANCHES           |
      | local, origin | main, alpha, gamma |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | gamma  | local, origin | gamma commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | gamma  | alpha  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | alpha  | git branch beta {{ sha 'WIP on beta' }} |
      |        | git checkout beta                       |
      | beta   | git reset {{ sha 'beta commit' }}       |
      |        | git push -u origin beta                 |
    And the current branch is now "beta"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branches and hierarchy exist
