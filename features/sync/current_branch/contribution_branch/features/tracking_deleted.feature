Feature: remove the contribution branch as soon as the tracking branch is gone, even if it has unpushed commits

  Background:
    Given the current branch is a contribution branch "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE      | FILE NAME  |
      | main         | local, origin | main commit  | main_file  |
      | contribution | local         | local commit | local_file |
    And origin deletes the "contribution" branch
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                    |
      | contribution | git fetch --prune --tags   |
      |              | git checkout main          |
      | main         | git branch -D contribution |
    And the current branch is now "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
    And it prints:
      """
      deleted branch "contribution"
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                     |
      | main   | git branch contribution {{ sha-before-run 'local commit' }} |
      |        | git checkout contribution                                   |
    And the current branch is now "contribution"
    And the initial commits exist
    And the initial branches and lineage exist
