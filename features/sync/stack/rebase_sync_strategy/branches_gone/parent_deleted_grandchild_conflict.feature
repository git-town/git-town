Feature: a grandchild branch has conflicts while its parent was deleted remotely

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                     | FILE NAME        | FILE CONTENT                             |
      | main     | local    | conflicting main commit     | conflicting_file | line 1\nline 2\nline 3: main content     |
      | branch-1 | local    | branch-1 commit             | child_file       | line 1: branch-1 content\nline 2\nline 3 |
      | branch-2 | local    | conflicting branch-2 commit | conflicting_file | line 1\nline 2: branch-2 content\nline 3 |
    And origin deletes the "branch-1" branch
    And the current branch is "branch-1" and the previous branch is "branch-2"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-1 | git fetch --prune --tags                                   |
      |          | git checkout main                                          |
      | main     | git -c rebase.updateRefs=false rebase origin/main          |
      |          | git push                                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
      |          | git checkout --theirs conflicting_file                     |
      |          | git add conflicting_file                                   |
      |          | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease                                |
      |          | git branch -D branch-1                                     |
      |          | git push --tags                                            |
    And no rebase is now in progress
    And all branches are now synchronized
    # TODO: the commits below are wrong.
    # branch-2 deletes the changes made by "conflicting main commit"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE                     | FILE NAME        | FILE CONTENT                             |
      | main     | local, origin | conflicting main commit     | conflicting_file | line 1\nline 2\nline 3: main content     |
      | branch-2 | local, origin | conflicting branch-2 commit | conflicting_file | line 1\nline 2: branch-2 content\nline 3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                |
      | branch-2 | git reset --hard {{ sha 'conflicting branch-2 commit' }}               |
      |          | git push --force-with-lease origin {{ sha 'initial commit' }}:branch-2 |
      |          | git branch branch-1 {{ sha 'branch-1 commit' }}                        |
      |          | git checkout branch-1                                                  |
    And the initial branches and lineage exist now
