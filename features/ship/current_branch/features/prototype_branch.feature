@skipWindows
Feature: shipping a prototype branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | prototype | local, origin | prototype commit |
    When I run "git-town ship" and enter "prototype done" for the commit message

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                           |
      | prototype | git fetch --prune --tags          |
      |           | git checkout main                 |
      | main      | git merge --squash --ff prototype |
      |           | git commit                        |
      |           | git push                          |
      |           | git push origin :prototype        |
      |           | git branch -D prototype           |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | prototype done |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git revert {{ sha 'prototype done' }}             |
      |        | git push                                          |
      |        | git branch prototype {{ sha 'prototype commit' }} |
      |        | git push -u origin prototype                      |
      |        | git checkout prototype                            |
    And the current branch is now "prototype"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                 |
      | main      | local, origin | prototype done          |
      |           |               | Revert "prototype done" |
      | prototype | local, origin | prototype commit        |
    And the initial branches and lineage exist
    And branch "prototype" is now prototype
