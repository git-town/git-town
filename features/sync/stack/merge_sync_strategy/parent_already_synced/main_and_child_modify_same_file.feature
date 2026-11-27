Feature: stop on a real conflict when the parent already merged a conflicting change from main

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT   |
      | main   | local, origin | main commit | file      | line 1\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME   | FILE CONTENT     |
      | branch-1 | local, origin | branch-1 commit | branch_file | branch-1 content |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                       |
      | branch-2 | local, origin | branch-2 commit | file      | line 1\nline 2 changed by branch-2 |
      | main     | local, origin | main update     | file      | line 1\nline 2 changed by main     |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And I ran "git-town sync" on branch "branch-1"
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                           |
      | branch-2 | git fetch --prune --tags          |
      |          | git checkout branch-1             |
      | branch-1 | git checkout branch-2             |
      | branch-2 | git merge --no-edit --ff branch-1 |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And a merge is now in progress
    And file "file" now has content:
      """
      line 1
      <<<<<<< HEAD
      line 2 changed by branch-2
      =======
      line 2 changed by main
      >>>>>>> branch-1
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND           |
      | branch-2 | git merge --abort |
    And no merge is now in progress
