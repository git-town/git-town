Feature: syncing a branch with independent changes where a commit was amended

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT               |
      | main   | local, origin | main commit | file      | line 0\n\nline 1\n\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                   |
      | branch-1 | local, origin | branch-1 commit A | file      | line 0\n\nline 1: branch-1 content A\n\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                       |
      | branch-2 | local, origin | branch-2 commit A | file      | line 0\n\nline 1: branch-1 content A\n\nline 2: branch-2 content A |
    And the current branch is "branch-2"
    And I ran "git-town sync"
    And I amend this commit
      | BRANCH   | LOCATION | MESSAGE           | FILE NAME | FILE CONTENT                                   |
      | branch-1 | local    | branch-1 commit B | file      | line 0\n\nline 1: branch-1 content B\n\nline 2 |
    And I amend this commit
      | BRANCH   | LOCATION | MESSAGE           | FILE NAME | FILE CONTENT                                                       |
      | branch-2 | local    | branch-2 commit B | file      | line 0\n\nline 1: branch-1 content A\n\nline 2: branch-2 content B |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                       |
      | main     | local, origin | main commit       | file      | line 0\n\nline 1\n\nline 2                                         |
      | branch-1 | local         | branch-1 commit B | file      | line 0\n\nline 1: branch-1 content B\n\nline 2                     |
      |          | origin        | branch-1 commit A | file      | line 0\n\nline 1: branch-1 content A\n\nline 2                     |
      | branch-2 | local         | branch-1 commit A | file      | line 0\n\nline 1: branch-1 content A\n\nline 2                     |
      |          |               | branch-2 commit B | file      | line 0\n\nline 1: branch-1 content A\n\nline 2: branch-2 content B |
      |          | origin        | branch-2 commit A | file      | line 0\n\nline 1: branch-1 content A\n\nline 2: branch-2 content A |
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
      |          | git push --force-with-lease --force-if-includes                                     |
    And no rebase is now in progress
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT                                                       |
      | main     | local, origin | main commit       | file      | line 0\n\nline 1\n\nline 2                                         |
      | branch-1 | local, origin | branch-1 commit B | file      | line 0\n\nline 1: branch-1 content B\n\nline 2                     |
      | branch-2 | local, origin | branch-2 commit B | file      | line 0\n\nline 1: branch-1 content B\n\nline 2: branch-2 content B |
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                   |
      | branch-2 | git reset --hard {{ sha 'branch-2 commit B' }}                            |
      |          | git push --force-with-lease origin {{ sha 'branch-2 commit A' }}:branch-2 |
      |          | git push --force-with-lease origin {{ sha 'branch-1 commit A' }}:branch-1 |
    And the initial commits exist now
    And the initial branches and lineage exist now
