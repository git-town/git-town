@skipWindows
Feature: shipping a prototype branch using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | prototype | local, origin | prototype commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "prototype"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                               |
      | prototype | git fetch --prune --tags              |
      |           | git checkout main                     |
      | main      | git merge --no-ff --edit -- prototype |
      |           | git push                              |
      |           | git push origin :prototype            |
      |           | git branch -D prototype               |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                  |
      | main   | local, origin | prototype commit         |
      |        |               | Merge branch 'prototype' |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git branch prototype {{ sha 'prototype commit' }} |
      |        | git push -u origin prototype                      |
      |        | git checkout prototype                            |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                  |
      | main   | local, origin | prototype commit         |
      |        |               | Merge branch 'prototype' |
    And branch "prototype" now has type "prototype"
