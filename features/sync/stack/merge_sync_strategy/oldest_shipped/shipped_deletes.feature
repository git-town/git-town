Feature: auto-resolve phantom merge conflicts

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT     |
      | main   | local, origin | main commit | file      | branch-1 content |
    And I ran "git-town hack branch-1"
    And I ran "git rm file"
    And I ran "git commit -m branch-1-commit"
    And I ran "git push -u origin branch-1"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME | FILE CONTENT     |
      | branch-2 | local    | branch-2 commit | file      | branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
      |          | git branch -D branch-1                            |
      |          | git checkout branch-2                             |
      | branch-2 | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (modify/delete): file deleted in main and modified in HEAD.  Version HEAD of file left in tree.
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                 |
      | branch-2 | git merge --abort                                       |
      |          | git checkout main                                       |
      | main     | git reset --hard {{ sha 'main commit' }}                |
      |          | git branch branch-1 {{ sha-initial 'branch-1-commit' }} |
      |          | git checkout branch-2                                   |
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
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME | FILE CONTENT     |
      | main     | local, origin | main commit                       | file      | branch-1 content |
      |          |               | branch-1-commit                   | file      | (deleted)        |
      | branch-2 | local, origin | branch-1-commit                   | file      | (deleted)        |
      |          |               | branch-2 commit                   | file      | branch-2 content |
      |          |               | Merge branch 'main' into branch-2 |           |                  |

  Scenario: resolve, commit, and continue
    When I ran "git add file"
    When I ran "git commit --no-edit"
    And I ran "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git merge --no-edit --ff origin/branch-2 |
      |          | git push                                 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                           | FILE NAME | FILE CONTENT     |
      | main     | local, origin | main commit                       | file      | branch-1 content |
      |          |               | branch-1-commit                   | file      | (deleted)        |
      | branch-2 | local, origin | branch-1-commit                   | file      | (deleted)        |
      |          |               | branch-2 commit                   | file      | branch-2 content |
      |          |               | Merge branch 'main' into branch-2 |           |                  |
