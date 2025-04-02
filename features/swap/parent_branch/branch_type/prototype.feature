Feature: swapping a branch with its prototype parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT    | LOCATIONS     |
      | prototype | prototype | main      | local, origin |
      | feature   | feature   | prototype | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | prototype | local, origin | prototype commit |
      | feature   | local, origin | feature commit   |
    And the current branch is "feature"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature   | git fetch --prune --tags                        |
      |           | git rebase --onto main prototype                |
      |           | git checkout prototype                          |
      | prototype | git rebase --onto feature main                  |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout feature                            |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature   | local, origin | feature commit   |
      | prototype | local, origin | prototype commit |
    And this lineage exists now
      | BRANCH    | PARENT  |
      | feature   | main    |
      | prototype | feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature   | git checkout prototype                          |
      | prototype | git reset --hard {{ sha 'prototype commit' }}   |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout feature                            |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial lineage exists now
