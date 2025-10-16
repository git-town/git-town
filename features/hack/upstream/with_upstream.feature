Feature: on a forked repo

  Background:
    Given a Git repo with origin
    And an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And the current branch is "main"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | main   | git fetch --prune --tags                            |
      |        | git fetch upstream main                             |
      |        | git -c rebase.updateRefs=false rebase upstream/main |
      |        | git push                                            |
      |        | git checkout -b new                                 |
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -D new |
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |
