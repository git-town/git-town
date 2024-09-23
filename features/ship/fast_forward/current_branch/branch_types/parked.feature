Feature: shipping a parked branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parked | local, origin | parked commit |
    And the current branch is "parked"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                    |
      | parked | git fetch --prune --tags   |
      |        | git checkout main          |
      | main   | git merge --ff-only parked |
      |        | git push                   |
      |        | git push origin :parked    |
      |        | git branch -D parked       |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parked commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git branch parked {{ sha 'parked commit' }} |
      |        | git push -u origin parked                   |
      |        | git checkout parked                         |
    And the current branch is now "parked"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parked commit |
      | parked | local, origin | parked commit |
    And the initial branches and lineage exist now
    And branch "parked" is now parked
