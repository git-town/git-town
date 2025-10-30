Feature: compress a branch with a merge commit

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main    | local, origin | main commit | main_file | main content |
      | feature | local, origin | commit 1    | file_1    | content 1    |
    And the current branch is "feature"
    And I ran "git merge main -m merge"
    And I ran "git push"
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git reset --soft main --                        |
      |         | git commit -m "commit 1"                        |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main    | local, origin | main commit | main_file | main content |
      | feature | local, origin | commit 1    | file_1    | content 1    |
    And file "file_1" still has content "content 1"
    And file "main_file" still has content "main content"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git reset --hard {{ sha 'merge' }}              |
      |         | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
