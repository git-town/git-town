Feature: remove a branch and all its children from a stack with dependent changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 1\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                             |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                               |
      | branch-2 | local, origin | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                 |
      | branch-3 | local, origin | branch-3 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes\nline 3: branch-3 changes |
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-2"
    When I run "git-town set-parent --none"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "branch-2" is now perennial
      """
    And this lineage exists now
      """
      branch-2
        branch-3

      main
        branch-1
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | local, origin | main commit     |
      | branch-2 | local, origin | main commit     |
      |          |               | branch-1 commit |
      |          |               | branch-2 commit |
      | branch-3 | local, origin | branch-3 commit |
      | branch-1 | local, origin | branch-1 commit |

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
