Feature: delete a branch within a branch chain

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "beta" and the previous branch is "alpha"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
      |        | git push origin :beta    |
      |        | git checkout alpha       |
      | alpha  | git branch -D beta       |
    And Git Town prints:
      """
      branch "gamma" is now a child of "alpha"
      """
    And the current branch is now "alpha"
    And no uncommitted files exist now
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
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | alpha  | git branch beta {{ sha 'beta commit' }} |
      |        | git push -u origin beta                 |
      |        | git checkout beta                       |
    And the current branch is now "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now
