Feature: stage all open changes through the long flag

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    When I run "git-town commit --up --all -m up-commit"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                           |
      | branch-2 | git add -A                        |
      |          | git checkout branch-1             |
      | branch-1 | git commit -m up-commit           |
      |          | git checkout branch-2             |
      | branch-2 | git merge --no-edit --ff branch-1 |
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | branch-1 | local    | up-commit | changes   | my changes   |

  Scenario: undo
    When I run "git-town undo"
    Then the initial branches and lineage exist now
    And the uncommitted file is untracked
    And the uncommitted file still exists
    And the initial commits exist now

  Scenario: undo with a corrupted snapshot backup
    When the open changes snapshot backup is corrupted
    And I run "git-town undo"
    Then Git Town prints the error:
      """
      the saved open changes snapshot is incomplete
      """
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | branch-1 | local    | up-commit | changes   | my changes   |
