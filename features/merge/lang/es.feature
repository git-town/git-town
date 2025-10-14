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
      | BRANCH | COMMAND                                  |
      | beta   | git checkout alpha                       |
      | alpha  | git reset --hard {{ sha 'beta commit' }} |
      |        | git branch -D beta                       |
    And Git Town prints:
      """
      Eliminada la rama beta
      """
    And this lineage exists now
      """
      main
        alpha
      """
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | alpha  | local    | alpha commit |
      |        |          | beta commit  |

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | alpha  | git reset --hard {{ sha 'alpha commit' }} |
      |        | git branch beta {{ sha 'beta commit' }}   |
      |        | git checkout beta                         |
    And the initial lineage exists now
    And the initial commits exist now
