@skipWindows
Feature: shipping a prototype branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | prototype | local, origin | prototype commit |
    And the current branch is "prototype"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                       |
      | prototype | git fetch --prune --tags      |
      |           | git checkout main             |
      | main      | git merge --ff-only prototype |
      |           | git push                      |
      |           | git push origin :prototype    |
      |           | git branch -D prototype       |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          |
      | main   | local, origin | prototype commit |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git branch prototype {{ sha 'prototype commit' }} |
      |        | git push -u origin prototype                      |
      |        | git checkout prototype                            |
    And the current branch is now "prototype"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          |
      | main   | local, origin | prototype commit |
    And the initial branches and lineage exist now
    And branch "prototype" is now prototype
