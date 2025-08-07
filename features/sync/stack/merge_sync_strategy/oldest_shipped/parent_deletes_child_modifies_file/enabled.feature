Feature: auto-resolve phantom merge conflicts in a synced stack where the parent deletes a file and gets shipped, and the child modifies the same file

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | local, origin | create file | file      | main content |
    And I ran "git-town hack branch-1"
    And I ran "git rm file"
    And I ran "git commit -m delete-file"
    And I ran "git push -u origin branch-1"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE     | FILE NAME | FILE CONTENT     |
      | branch-2 | local    | change file | file      | branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
      |          | git branch -D branch-1                            |
      |          | git checkout branch-2                             |
      | branch-2 | git merge --no-edit --ff main                     |
    # TODO: auto-resolve this phantom merge conflict.
    # Branch-1 deletes "file" and branch-2 creates it again.
    # Branch-2 was properly synced with branch-1.
    # When branch-1 got shipped, and the user syncs, they shouldn't need to tell Git again that branch-2 should re-create the "file".
    #
    # note: this uses the merge sync strategy, not rebase
    #
    # This could be detected the same way we detect phantom merge conflicts:
    # If one variant of the conflict matches the root branch, pick the other variant?
    And Git Town prints the error:
      """
      CONFLICT (modify/delete): file deleted in main and modified in HEAD.  Version HEAD of file left in tree.
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                             |
      | branch-2 | git merge --abort                                   |
      |          | git checkout main                                   |
      | main     | git reset --hard {{ sha 'create file' }}            |
      |          | git branch branch-1 {{ sha-initial 'delete-file' }} |
      |          | git checkout branch-2                               |
    And the initial commits exist now
    And no merge is now in progress

  Scenario: continue with unresolved conflict
    When I run "git town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is now in progress

  Scenario: resolve and continue
    When I ran "git add file"
    And I ran "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git commit --no-edit                     |
      |          | git merge --no-edit --ff origin/branch-2 |
      |          | git push                                 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME      | FILE CONTENT     |
      | main     | local, origin | create file                       | file           | main content     |
      |          |               | delete-file                       | file (deleted) |                  |
      | branch-2 | local, origin | delete-file                       | file (deleted) |                  |
      |          |               | change file                       | file           | branch-2 content |
      |          |               | Merge branch 'main' into branch-2 |                |                  |

  Scenario: resolve, commit, and continue
    When I ran "git add file"
    And I ran "git commit --no-edit"
    And I ran "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git merge --no-edit --ff origin/branch-2 |
      |          | git push                                 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME      | FILE CONTENT     |
      | main     | local, origin | create file                       | file           | main content     |
      |          |               | delete-file                       | file (deleted) |                  |
      | branch-2 | local, origin | delete-file                       | file (deleted) |                  |
      |          |               | change file                       | file           | branch-2 content |
      |          |               | Merge branch 'main' into branch-2 |                |                  |
