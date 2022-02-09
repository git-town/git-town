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
    And I am on the "gamma" branch
    And my workspace has an uncommitted file
    When I run "git-town kill beta"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | gamma  | git fetch --prune --tags |
      |        | git push origin :beta    |
      |        | git branch -D beta       |
    And I am now on the "gamma" branch
    And my workspace still contains my uncommitted file
    And the branches are now
      | REPOSITORY    | BRANCHES           |
      | local, origin | main, alpha, gamma |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | gamma  | local, origin | gamma commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | alpha  | main   |
      | gamma  | alpha  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | gamma  | git branch beta {{ sha 'beta commit' }} |
      |        | git push -u origin beta                 |
    And I am now on the "gamma" branch
    And my workspace has the uncommitted file again
    And now the initial commits exist
    And the initial branch setup exists
