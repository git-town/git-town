Feature: merging in Spanish

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | local    | alpha commit |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS |
      | beta | feature | alpha  | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | beta   | local    | beta commit |
    And the current branch is "beta"
    When I run "git-town merge" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | beta   | git branch -D alpha |
    And Git Town prints:
      """
      Eliminada la rama alpha
      """
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | beta   | local    | alpha commit |
      |        |          | beta commit  |

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git branch alpha {{ sha-before-run 'alpha commit' }} |
    And the initial commits exist now
    And the initial lineage exists now
