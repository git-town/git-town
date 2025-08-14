Feature: syncing a branch with dependent changes where a commit was amended

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT           |
      | main   | local, origin | main commit | file      | line 0\nline 1\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                               |
      | branch-1 | local, origin | branch-1 commit A | file      | line 0\nline 1: branch-1 content A\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                   |
      | branch-2 | local, origin | branch-2 commit A | file      | line 0\nline 1: branch-1 content A\nline 2: branch-2 content A |
    And the current branch is "branch-2"
    And I ran "git-town sync"
    And I amend this commit
      | BRANCH   | LOCATION | MESSAGE           | FILE NAME | FILE CONTENT                               |
      | branch-1 | local    | branch-1 commit B | file      | line 0\nline 1: branch-1 content B\nline 2 |
    And I amend this commit
      | BRANCH   | LOCATION | MESSAGE           | FILE NAME | FILE CONTENT                                                   |
      | branch-2 | local    | branch-2 commit B | file      | line 0\nline 1: branch-1 content A\nline 2: branch-2 content B |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                   |
      | main     | local, origin | main commit       | file      | line 0\nline 1\nline 2                                         |
      | branch-1 | local         | branch-1 commit B | file      | line 0\nline 1: branch-1 content B\nline 2                     |
      |          | origin        | branch-1 commit A | file      | line 0\nline 1: branch-1 content A\nline 2                     |
      | branch-2 | local         | branch-1 commit A | file      | line 0\nline 1: branch-1 content A\nline 2                     |
      |          |               | branch-2 commit B | file      | line 0\nline 1: branch-1 content A\nline 2: branch-2 content B |
      |          | origin        | branch-2 commit A | file      | line 0\nline 1: branch-1 content A\nline 2: branch-2 content A |
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                             |
      | branch-2 | git fetch --prune --tags                                                            |
      |          | git checkout branch-1                                                               |
      | branch-1 | git push --force-with-lease --force-if-includes                                     |
      |          | git checkout branch-2                                                               |
      | branch-2 | git push --force-with-lease --force-if-includes                                     |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha 'branch-1 commit A' }} |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And file "file" now has content:
      """
      line 0
      <<<<<<< HEAD
      line 1: branch-1 content B
      line 2
      =======
      line 1: branch-1 content A
      line 2: branch-2 content B
      >>>>>>> {{ sha-short 'branch-2 commit B' }} (branch-2 commit B)
      """
    When I resolve the conflict in "file" with:
      """
      line 0
      line 1: branch-1 content B
      line 2: branch-2 content B
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                   |
      | main     | local, origin | main commit       | file      | line 0\nline 1\nline 2                                         |
      | branch-1 | local, origin | branch-1 commit B | file      | line 0\nline 1: branch-1 content B\nline 2                     |
      | branch-2 | local, origin | branch-2 commit B | file      | line 0\nline 1: branch-1 content B\nline 2: branch-2 content B |
