Feature: handle conflicts between the supplied feature branch and the main branch

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And the current branch is "other"
    And an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                         |
      | other  | git fetch --prune --tags        |
      |        | git add -A                      |
      |        | git stash                       |
      |        | git checkout main               |
      | main   | git merge --squash --ff feature |
      |        | git reset --hard                |
      |        | git checkout other              |
      | other  | git stash pop                   |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And no merge is in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is now "other"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial branches and lineage exist
