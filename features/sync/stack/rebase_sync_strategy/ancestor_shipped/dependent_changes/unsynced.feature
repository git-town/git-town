Feature: let the user resolve the merge conflict in an unsynced stack where the parent modifies a file and gets shipped, and the child modifies the same file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                         |
      | branch-1 | local, origin | branch-1 commit | file      | line 1 changed by branch-1\n\nline 2 |
      | branch-2 | local         | branch-2 commit | file      | line 1\n\nline 2 changed by branch-2 |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-2 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1 changed by branch-1

      line 2
      =======
      line 1

      line 2 changed by branch-2
      >>>>>>> {{ sha-short 'branch-2 commit' }} (branch-2 commit)
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                     |
      | branch-2 | git rebase --abort                          |
      |          | git checkout main                           |
      | main     | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout branch-2                       |
    And the initial branches and lineage exist now

  Scenario: resolve and continue
    When I resolve the conflict in "file" with:
      """
      line 1 changed by branch-1

      line 2 changed by branch-2
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                               |
      | branch-2 | GIT_EDITOR=true git rebase --continue |
      |          | git push --force-with-lease           |
      |          | git branch -D branch-1                |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                             |
      | main     | local, origin | branch-1 commit | file      | line 1 changed by branch-1\n\nline 2                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 1 changed by branch-1\n\nline 2 changed by branch-2 |
