Feature: delete a branch within a branch chain

  Background:
    Given my repo has a feature branch "alpha"
    And my repo has a feature branch "beta" as a child of "alpha"
    And my repo has a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, remote | alpha commit |
      | beta   | local, remote | beta commit  |
      | gamma  | local, remote | gamma commit |
    And I am on the "beta" branch
    And my workspace has an uncommitted file
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
    And I am now on the "alpha" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY    | BRANCHES           |
      | local, remote | main, alpha, gamma |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, remote | alpha commit |
      | gamma  | local, remote | gamma commit |
    And Git Town is now aware of this branch hierarchy
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
    And I am now on the "beta" branch
    And my workspace has the uncommitted file again
    And now the initial commits exist
    And my repo now has its initial branches and branch hierarchy
