Feature: hack in Spanish

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | origin   | main commit |
    And the current branch is "main"
    When I run "git-town hack new" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout -b new                               |
    And Git Town prints:
      """
      Cambiado a nueva rama 'new'
      """
    And this lineage exists now
      """
      main
        new
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | new    | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch -D new                           |
    And Git Town prints:
      """
      Eliminada la rama new
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
