Feature: remove a branch from a stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1 | file_1    | content 1    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-2 | local, origin | commit 2 | file_2    | content 2    |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-3 | local, origin | commit 3 | file_1    | content 3    |
    And the current branch is "branch-3"
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town set-parent" and enter into the dialog:
      | DIALOG                 | KEYS            |
      | parent branch of child | down down enter |

  Scenario: result
    Then Git Town prints:
      """
      Selected parent branch for "branch-3": main
      """
    And Git Town prints:
      """
      branch "branch-3" is now a child of "main"
      """
    And Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | git pull                                        |
      |          | git rebase --onto main branch-2 branch-3        |
      |          | git rm file_1                                   |
      |          | git -c core.editor=true rebase --continue       |
      |          | git push --force-with-lease --force-if-includes |
    And the current branch is still "branch-3"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1 | file_1    | content 1    |
      | branch-2 | local, origin | commit 2 | file_2    | content 2    |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | branch-1 |
      | branch-3 | main     |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}           |
      |          | git push --force-with-lease --force-if-includes |
    And the current branch is still "branch-3"
    And the initial commits exist now
    And the initial branches and lineage exist now
