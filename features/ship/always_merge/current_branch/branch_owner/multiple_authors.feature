Feature: ship a coworker's feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE            | AUTHOR                            |
      | feature | local, origin | developer commit 1 | developer <developer@example.com> |
      |         |               | developer commit 2 | developer <developer@example.com> |
      |         |               | coworker commit    | coworker <coworker@example.com>   |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "feature"
    When I run "git-town ship" and close the editor

  Scenario: result
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                | AUTHOR                            |
      | main   | local, origin | developer commit 1     | developer <developer@example.com> |
      |        |               | developer commit 2     | developer <developer@example.com> |
      |        |               | coworker commit        | coworker <coworker@example.com>   |
      |        |               | Merge branch 'feature' | user <email@example.com>          |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature {{ sha 'coworker commit' }} |
      |        | git push -u origin feature                     |
      |        | git checkout feature                           |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | developer commit 1     |
      |        |               | developer commit 2     |
      |        |               | coworker commit        |
      |        |               | Merge branch 'feature' |
