Feature: compressing a branch when its parent received additional commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          | FILE NAME    | FILE CONTENT      |
      | feature | local, origin | feature commit 1 | feature_file | feature content 1 |
      | feature | local, origin | feature commit 2 | feature_file | feature content 2 |
      | main    | local, origin | main commit      | main_file    | main content      |
    And the current branch is "feature"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" is not in sync with its parent, please run "git town sync" and try again
      """
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
