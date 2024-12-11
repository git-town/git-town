@messyoutput
Feature: remove a branch from a stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME |
      | branch-1 | local, origin | commit 1 | file_1    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME |
      | branch-2 | local, origin | commit 2 | file_2    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME |
      | branch-3 | local, origin | commit 3 | file_3    |
    And the current branch is "branch-3"
    And local Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS                 |
      | parent branch of child | down down down enter |

  Scenario: result
    Then Git Town prints:
      """
      Selected parent branch for "branch-3": branch-1
      """
    And Git Town prints:
      """
      branch "branch-3" is now a child of "branch-1"
      """
    And Git Town runs no commands
    And the current branch is still "branch-3"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
      | branch-3 | local, origin | commit 2 |
      |          |               | commit 3 |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | branch-1 |
      | branch-3 | branch-1 |
    And the branches contain these files:
      | BRANCH   | NAME   |
      | branch-1 | file_1 |
      | branch-2 | file_1 |
      |          | file_2 |
      | branch-3 | file_1 |
      |          | file_2 |
      |          | file_3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "branch-3"
    And the initial commits exist now
    And the initial branches and lineage exist now
