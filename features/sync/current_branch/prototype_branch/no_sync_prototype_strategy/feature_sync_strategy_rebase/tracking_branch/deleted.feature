Feature: remove a prototype branch as soon as its tracking branch is gone, even if it has unpushed commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE      | FILE NAME  |
      | main      | local, origin | main commit  | main_file  |
      | prototype | local         | local commit | local_file |
    And the current branch is "prototype"
    And origin deletes the "prototype" branch
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                          |
      | prototype | git fetch --prune --tags         |
      |           | git checkout main                |
      | main      | git rebase --onto main prototype |
      |           | git branch -D prototype          |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And Git Town prints:
      """
      deleted branch "prototype"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                  |
      | main   | git branch prototype {{ sha-before-run 'local commit' }} |
      |        | git checkout prototype                                   |
    And the initial commits exist now
    And the initial branches and lineage exist now
