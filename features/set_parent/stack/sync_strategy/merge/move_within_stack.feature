@messyoutput
Feature: remove a branch from a stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME |
      | branch-1 | local, origin | existing commit | file_1    |
      | branch-2 | local, origin | existing commit | file_2    |
      | branch-3 | local, origin | existing commit | file_3    |
    And the current branch is "branch-2"
    And local Git Town setting "sync-feature-strategy" is "merge"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS       |
      | parent branch of child | down enter |

  @this
  Scenario: result
    Then Git Town prints:
      """
      Selected parent branch for "branch-2": main
      """
    And Git Town runs no commands
    And the current branch is still "branch-2"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | main     |
      | branch-3 | branch-2 |
    And the branches contain these files:
      | BRANCH   | NAME   |
      | branch-1 | file_1 |
      | branch-2 | file_2 |
      | branch-3 | file_2 |
      |          | file_3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial branches and lineage exist now
