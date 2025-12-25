Feature: prepend a branch in Spanish

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And the current branch is "old"
    When I run "git-town prepend parent" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git fetch --prune --tags    |
      |        | git checkout -b parent main |
    And Git Town prints:
      """
      Cambiado a nueva rama 'parent'
      """
    And this lineage exists now
      """
      main
        parent
          old
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And Git Town prints:
      """
      Eliminada la rama parent
      """
    And the initial lineage exists now
    And the initial commits exist now
