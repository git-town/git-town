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
    And the current branch is "branch-2"
    And local Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS     |
      | parent branch of child | up enter |

  @this
  Scenario: result
    Then Git Town prints:
      """
      Selected parent branch for "branch-2": main
      """
    And Git Town runs the commands
      | BRANCH   | COMMAND                                      |
      | branch-2 | git pull                                     |
      |          | git rebase --onto main branch-1 branch-2     |
      |          | git push --force-with-lease                  |
      |          | git checkout branch-3                        |
      | branch-3 | git pull                                     |
      |          | git rebase --onto branch-2 branch-2 branch-3 |
      |          | git push --force-with-lease                  |
      |          | git checkout branch-2                        |
    # And inspect the repo
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
      | branch-2 | local, origin | commit 2 |
      | branch-3 | local, origin | commit 1 |
      |          |               | commit 3 |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | main     |
      | branch-3 | branch-2 |
    And the branches contain these files:
      | BRANCH   | NAME   |
      | branch-1 | file_1 |
      | branch-2 | file_2 |
      | branch-3 | file_3 |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial branches and lineage exist now
