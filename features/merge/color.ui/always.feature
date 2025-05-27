Feature: merging a branch in a stack that is fully in sync

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | beta   | local, origin | beta commit |
    And the current branch is "beta"
    And local Git setting "color.ui" is "always"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git fetch --prune --tags                        |
      |        | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'beta commit' }}        |
      |        | git branch -D beta                              |
      |        | git push --force-with-lease --force-if-includes |
      |        | git push origin :beta                           |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      |        |               | beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | alpha  | git reset --hard {{ sha-initial 'alpha commit' }} |
      |        | git push --force-with-lease --force-if-includes   |
      |        | git branch beta {{ sha 'beta commit' }}           |
      |        | git push -u origin beta                           |
      |        | git checkout beta                                 |
    And the initial commits exist now
    And the initial lineage exists now
