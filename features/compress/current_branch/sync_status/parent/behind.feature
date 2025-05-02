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
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git reset --soft main                           |
      |         | git commit -m "feature commit 1"                |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE          |
      | main    | local, origin | main commit      |
      | feature | local, origin | feature commit 1 |
    And the branches contain these files:
      | BRANCH  | NAME         |
      | feature | feature_file |
      | main    | main_file    |
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git reset --hard {{ sha 'feature commit 2' }}   |
      |         | git push --force-with-lease --force-if-includes |
    And the initial commits exist now
    And the initial branches and lineage exist now
