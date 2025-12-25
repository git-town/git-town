Feature: on the main branch with an upstream repo

  Background:
    Given a Git repo with origin
    And an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And Git setting "git-town.sync-upstream" is "true"
    And the current branch is "main"
    And I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | main   | git fetch --prune --tags                            |
      |        | git fetch upstream main                             |
      |        | git -c rebase.updateRefs=false rebase upstream/main |
      |        | git push                                            |
      |        | git push --tags                                     |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |
