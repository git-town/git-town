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
    And an uncommitted file
    And the current branch is "prototype" and the previous branch is "previous"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                                          |
      | prototype | git fetch --prune --tags                         |
      |           | git push origin :prototype                       |
      |           | git add -A                                       |
      |           | git commit -m "Committing WIP for git town undo" |
      |           | git checkout previous                            |
      | previous  | git branch -D prototype                          |
    And the current branch is now "previous"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, previous |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | previous | local, origin | previous commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | previous | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                                           |
      | previous  | git push origin {{ sha 'prototype commit' }}:refs/heads/prototype |
      |           | git branch prototype {{ sha 'Committing WIP for git town undo' }} |
      |           | git checkout prototype                                            |
      | prototype | git reset --soft HEAD~1                                           |
    And the current branch is now "prototype"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
    And branch "prototype" is now prototype
