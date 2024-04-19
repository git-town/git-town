Feature: handle conflicts between the shipped branch and the main branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git reset --hard                |
      |         | git checkout feature            |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And the current branch is still "feature"
    And no merge is in progress

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
    And it runs no commands
    And the current branch is still "feature"
    And no merge is in progress
    And the initial commits exist
    And the initial branches and lineage exist
