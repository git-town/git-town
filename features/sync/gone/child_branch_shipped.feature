Feature: sync a child branch that was deleted at the remote

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
    And origin deletes the "beta" branch
    And the current branch is "alpha"
    When I run "git-town sync --gone"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
      |        | git branch -D beta       |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | gamma  | local, origin | gamma commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git branch beta {{ sha-initial 'beta commit' }} |
    And the initial lineage exists now
    And the initial commits exist now
