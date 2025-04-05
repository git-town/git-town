Feature: merging when the branch is not in sync with its tracking branch

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
      | BRANCH | LOCATION | MESSAGE            |
      | beta   | origin   | remote beta commit |
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "beta" is not in sync with its tracking branch, please run "git town sync" and try again
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial lineage exists now
