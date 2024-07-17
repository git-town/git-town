Feature: abort the ship via empty commit message

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | main commit    | main_file        | main content    |
      | feature | local         | feature commit | conflicting_file | feature content |
    And the current branch is "other"
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town ship feature" and enter an empty commit message

  @skipWindows
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                         |
      | other  | git fetch --prune --tags        |
      |        | git add -A                      |
      |        | git stash                       |
      |        | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git commit                      |
      |        | git reset --hard                |
      |        | git checkout other              |
      | other  | git stash pop                   |
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "other"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, origin | main commit    |
      | feature | local         | feature commit |
    And the initial lineage exists
