Feature: delete the current feature branch in Spanish

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And the current branch is "current" and the previous branch is "other"
    When I run "git-town delete" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
      |         | git push origin :current |
      |         | git checkout other       |
      | other   | git branch -D current    |
    And Git Town prints:
      """
      Eliminada la rama current
      """
    And this lineage exists now
      """
      main
        other
      """
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch current {{ sha 'current commit' }} |
      |        | git push -u origin current                    |
      |        | git checkout current                          |
    And Git Town prints:
      """
      Cambiado a rama 'current'
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
