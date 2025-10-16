Feature: an ancestor in a stack with dependent changes was deleted remotely

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT   |
      | main   | local, origin | main commit | file      | line 1\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                     |
      | branch-1 | local, origin | branch-1 commit | file      | line 1: branch-1 changes\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         | FILE NAME | FILE CONTENT                                       |
      | branch-2 | local    | branch-2 commit | file      | line 1: branch-1 changes\nline 2: branch-2 changes |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "branch-1" branch
    And the current branch is "branch-1" and the previous branch is "branch-2"
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-1 | git fetch --prune --tags                                   |
      |          | git checkout branch-2                                      |
      | branch-2 | git pull                                                   |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1 |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1
      line 2
      =======
      line 1: branch-1 changes
      line 2: branch-2 changes
      >>>>>>> {{ sha-short 'branch-2 commit' }} (branch-2 commit)
      """
    When I resolve the conflict in "file" with:
      """
      line 1
      line 2: branch-2 changes
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
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                     |
      | main     | local, origin | main commit     | file      | line 1\nline 2                   |
      | branch-2 | local, origin | branch-2 commit | file      | line 1\nline 2: branch-2 changes |
