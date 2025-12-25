Feature: delete the current prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
      | previous  | feature   | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | previous  | local, origin | previous commit  |
      | prototype | local, origin | prototype commit |
    And the current branch is "prototype" and the previous branch is "previous"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                    |
      | prototype | git fetch --prune --tags   |
      |           | git push origin :prototype |
      |           | git checkout previous      |
      | previous  | git branch -D prototype    |
    And this lineage exists now
      """
      main
        previous
      """
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, previous |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | previous | local, origin | previous commit |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | previous | git branch prototype {{ sha 'prototype commit' }} |
      |          | git push -u origin prototype                      |
      |          | git checkout prototype                            |
    And the initial branches and lineage exist now
    And branch "prototype" now has type "prototype"
    And the initial commits exist now
