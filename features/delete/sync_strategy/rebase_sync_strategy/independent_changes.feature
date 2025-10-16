Feature: deleting a branch from a stack with independent changes

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                       |
      | main   | local, origin | main commit | file      | line 0: main content\n\nline 1\n\nline 2\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                         |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                           |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2: branch-2 content\n\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                             |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2: branch-2 content\n\nline 3: branch-3 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-2"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git push origin :branch-2                                      |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
      |          | git push --force-with-lease                                    |
      |          | git branch -D branch-2                                         |
    And this lineage exists now
      """
      main
        branch-1
          branch-3
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                 |
      | local, origin | main, branch-1, branch-3 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                           |
      | main     | local, origin | main commit     | file      | line 0: main content\n\nline 1\n\nline 2\n\nline 3                                     |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\n\nline 1: branch-1 content\n\nline 2\n\nline 3: branch-3 content |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                 |
      | branch-3 | git reset --hard {{ sha-initial 'branch-3 commit' }}    |
      |          | git push --force-with-lease --force-if-includes         |
      |          | git branch branch-2 {{ sha-initial 'branch-2 commit' }} |
      |          | git push -u origin branch-2                             |
      |          | git checkout branch-2                                   |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES                           |
      | local, origin | main, branch-1, branch-2, branch-3 |
