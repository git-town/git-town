Feature: sync inside a folder that doesn't exist on the main branch

  Background:
    Given the feature branches "current" and "other"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | conflicting main commit    | conflicting_file | main content    |
      | current | local         | conflicting current commit | conflicting_file | current content |
      |         |               | folder commit              | new_folder/file1 |                 |
      | other   | local, origin | other commit               | file2            |                 |
    And the current branch is "current"
    And an uncommitted file
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | current | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout current               |
      | current | git merge --no-edit origin/current |
      |         | git merge --no-edit main           |
    And the current branch is still "current"
    And the uncommitted file is stashed
    And a merge is now in progress
    And it prints the error:
      """
      exit status 1
      """

  Scenario: undo
    When I run "git-town undo" in the "new_folder" folder
    Then it runs the commands
      | BRANCH  | COMMAND           |
      | current | git merge --abort |
      |         | git stash pop     |
    And the current branch is still "current"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue" in the "new_folder" folder
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "current"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" in the "new_folder" folder
    Then it runs the commands
      | BRANCH  | COMMAND                          |
      | current | git commit --no-edit             |
      |         | git push                         |
      |         | git checkout other               |
      | other   | git merge --no-edit origin/other |
      |         | git merge --no-edit main         |
      |         | git push                         |
      |         | git checkout current             |
      | current | git push --tags                  |
      |         | git stash pop                    |
    And all branches are now synchronized
    And the current branch is still "current"
    And the uncommitted file still exists
    And no merge is in progress
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | conflicting main commit          |
      | current | local, origin | conflicting current commit       |
      |         |               | folder commit                    |
      |         |               | conflicting main commit          |
      |         |               | Merge branch 'main' into current |
      | other   | local, origin | other commit                     |
      |         |               | conflicting main commit          |
      |         |               | Merge branch 'main' into other   |
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | current | conflicting_file | resolved content |
      |         | new_folder/file1 |                  |
      | other   | conflicting_file | main content     |
      |         | file2            |                  |
