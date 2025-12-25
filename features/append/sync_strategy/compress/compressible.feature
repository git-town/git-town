Feature: append a new feature branch in a clean workspace using the "compress" sync strategy with compressible commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      | existing | local, origin | existing commit 2 |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "existing"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                           |
      | existing | git fetch --prune --tags          |
      |          | git reset --soft main --          |
      |          | git commit -m "existing commit 1" |
      |          | git push --force-with-lease       |
      |          | git checkout -b new               |
    And this lineage exists now
      """
      main
        existing
          new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                |
      | new      | git checkout existing                                  |
      | existing | git reset --hard {{ sha-initial 'existing commit 2' }} |
      |          | git push --force-with-lease --force-if-includes        |
      |          | git branch -D new                                      |
    And the initial lineage exists now
    And the initial commits exist now
