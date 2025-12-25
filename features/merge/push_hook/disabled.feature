Feature: merging a branch with disabled push-hook

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
    And Git setting "git-town.push-hook" is "false"
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                     |
      | beta   | git fetch --prune --tags                                    |
      |        | git checkout alpha                                          |
      | alpha  | git reset --hard {{ sha 'beta commit' }}                    |
      |        | git push origin :beta                                       |
      |        | git branch -D beta                                          |
      |        | git push --force-with-lease --force-if-includes --no-verify |
    And this lineage exists now
      """
      main
        alpha
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      |        |               | beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                     |
      | alpha  | git reset --hard {{ sha 'alpha commit' }}                   |
      |        | git push --force-with-lease --force-if-includes --no-verify |
      |        | git branch beta {{ sha 'beta commit' }}                     |
      |        | git push --no-verify -u origin beta                         |
      |        | git checkout beta                                           |
    And the initial lineage exists now
    And the initial commits exist now
