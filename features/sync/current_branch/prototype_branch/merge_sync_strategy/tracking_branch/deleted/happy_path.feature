Feature: remove a prototype branch as soon as its tracking branch is gone, even if it has unpushed commits

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE      | FILE NAME  |
      | main      | local, origin | main commit  | main_file  |
      | prototype | local         | local commit | local_file |
    And origin deletes the "prototype" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | prototype | git fetch --prune --tags |
      |           | git checkout main        |
      | main      | git rebase origin/main   |
      |           | git branch -D prototype  |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And it prints:
      """
      deleted branch "prototype"
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                  |
      | main   | git branch prototype {{ sha-before-run 'local commit' }} |
      |        | git checkout prototype                                   |
    And the current branch is now "prototype"
    And the initial commits exist
    And the initial branches and lineage exist
