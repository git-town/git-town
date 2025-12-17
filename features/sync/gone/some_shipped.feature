Feature: sync only branches whose remote is gone

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "alpha"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    And origin ships the "beta" branch using the "squash-merge" ship-strategy
    When I run "git-town sync --gone"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | alpha  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git branch -D alpha                               |
      |        | git branch -D beta                                |
      |        | git checkout gamma                                |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | alpha commit |
      |        |               | beta commit  |
      | gamma  | local, origin | gamma commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | gamma  | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch alpha {{ sha 'alpha commit' }}       |
      |        | git branch beta {{ sha-initial 'beta commit' }} |
      |        | git checkout alpha                              |
    And the initial lineage exists now
    And the initial commits exist now
