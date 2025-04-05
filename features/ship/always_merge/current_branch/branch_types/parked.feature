Feature: shipping a parked branch using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parked | local, origin | parked commit |
    And the current branch is "parked"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                            |
      | parked | git fetch --prune --tags           |
      |        | git checkout main                  |
      | main   | git merge --no-ff --edit -- parked |
      |        | git push                           |
      |        | git push origin :parked            |
      |        | git branch -D parked               |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | main   | local, origin | parked commit         |
      |        |               | Merge branch 'parked' |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git branch parked {{ sha 'parked commit' }} |
      |        | git push -u origin parked                   |
      |        | git checkout parked                         |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | main   | local, origin | parked commit         |
      |        |               | Merge branch 'parked' |
    And the initial branches and lineage exist now
    And branch "parked" now has type "parked"
