Feature: merging a branch in offline mode

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
    And offline mode is enabled
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                  |
      | beta   | git checkout alpha                       |
      | alpha  | git reset --hard {{ sha 'beta commit' }} |
      |        | git branch -D beta                       |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      |        | local         | beta commit  |
      | beta   | origin        | alpha commit |
      |        |               | beta commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | alpha  | git reset --hard {{ sha 'alpha commit' }} |
      |        | git branch beta {{ sha 'beta commit' }}   |
      |        | git checkout beta                         |
    And the initial commits exist now
    And the initial lineage exists now
