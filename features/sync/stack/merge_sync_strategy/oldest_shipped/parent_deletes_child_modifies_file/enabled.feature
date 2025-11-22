Feature: auto-resolve phantom merge conflicts in a synced stack where the parent deletes a file and gets shipped, and the child modifies the same file

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "merge"
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
    And I ran "git-town sync"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
      |          | git checkout branch-2                             |
      | branch-2 | git branch -D branch-1                            |
      |          | git merge --no-edit --ff main                     |
    # TODO: auto-resolve this phantom merge conflict
    #
    # This requires storing the SHA of branches after they were synced the last time.
    #
    # If Git Town stores this information.
    # We know that at the beginning of the second sync call
    # branch-2 is still in sync with its parent (branch-1).
    # That's because the initial SHA of either branch is the same as
    # their SHA was at the end of the last sync
    # (when they were guaranteed to be in sync).
    #
    # In the conflict, branch-2 modifies the file while the main branch deletes it.
    # Since branch-1 (which is in sync) also deletes it,
    # it is safe to keep the version on branch-2.
    And Git Town prints the error:
      """
      CONFLICT (modify/delete): file deleted in main and modified in HEAD.
      """
    And a merge is now in progress
    And file "file" now has content:
      """
      branch-2 content
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                             |
      | branch-2 | git merge --abort                                   |
      |          | git checkout main                                   |
      | main     | git reset --hard {{ sha 'create file' }}            |
      |          | git branch branch-1 {{ sha-initial 'delete-file' }} |
      |          | git checkout branch-2                               |
    And no merge is now in progress
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
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
      | BRANCH   | COMMAND              |
      | branch-2 | git commit --no-edit |
      |          | git push             |
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
      | BRANCH   | COMMAND  |
      | branch-2 | git push |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME      | FILE CONTENT     |
      | main     | local, origin | create file                       | file           | main content     |
      |          |               | delete-file                       | file (deleted) |                  |
      | branch-2 | local, origin | delete-file                       | file (deleted) |                  |
      |          |               | change file                       | file           | branch-2 content |
      |          |               | Merge branch 'main' into branch-2 |                |                  |
