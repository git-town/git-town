Feature: sync inside a folder that doesn't exist on the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | current | local         | conflicting current commit | conflicting_file | current content |
      |         |               | folder commit              | new_folder/file1 |                 |
      | other   | local, origin | other commit               | file2            |                 |
    And the current branch is "current"
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                       |
      | current | git fetch --prune --tags      |
      |         | git merge --no-edit --ff main |
    And Git Town prints the error:
      """
      git merge conflict
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo" in the "new_folder" folder
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | current | git merge --abort |
    And no merge is now in progress
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue" in the "new_folder" folder
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" in the "new_folder" folder
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | current | git commit --no-edit                    |
      |         | git merge --no-edit --ff origin/current |
      |         | git push                                |
      |         | git checkout other                      |
      | other   | git merge --no-edit --ff main           |
      |         | git push                                |
      |         | git checkout current                    |
      | current | git push --tags                         |
    And no merge is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | conflicting main commit          |
      | current | local, origin | conflicting current commit       |
      |         |               | folder commit                    |
      |         |               | Merge branch 'main' into current |
      | other   | local, origin | other commit                     |
      |         |               | Merge branch 'main' into other   |
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | current | conflicting_file | resolved content |
      |         | new_folder/file1 |                  |
      | other   | conflicting_file | main content     |
      |         | file2            |                  |
