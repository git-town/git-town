Feature: deleting a branch that conflicts with the main branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                 |
      | main   | local, origin | main commit | file      | line 0: main content\nline 1\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                   |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                     |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                       |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3: branch-3 content |
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
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And Git Town prints the error:
      """
      git rebase conflict
      """
    # This seems wrong. It should not rebase branch-3 onto main,
    # it should rebase it onto branch-1.
    And file "file" now has content:
      """
      line 0: main content
      line 1: branch-1 content
      <<<<<<< HEAD
      line 2
      line 3
      =======
      line 2: branch-2 content
      line 3: branch-3 content
      >>>>>>> {{ sha-short 'branch-3 commit' }} (branch-3 commit)
      """
    And a rebase is now in progress

  @this
  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                 |
      | branch-3 | git rebase --abort                                                      |
      |          | git push origin {{ sha-initial 'branch-2 commit' }}:refs/heads/branch-2 |
      |          | git checkout branch-2                                                   |
    And the branches are now
      | REPOSITORY    | BRANCHES                           |
      | local, origin | main, branch-1, branch-2, branch-3 |
    And the initial lineage exists now

  Scenario:
    And the branches are now
      | REPOSITORY    | BRANCHES                 |
      | local, origin | main, branch-1, branch-3 |
    # TODO: the commits below are wrong.
    # Branch-3 still contains the changes from branch-2.
    # These changes should have been removed when branch-2 was deleted.
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                      |
      | main     | local, origin | main commit     | file      | line 0: main content\nline 1\nline2\nline 3                                                       |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\nline 1: branch-1 content\nline2\n\nline 3                                   |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\nline 1: branch-1 content\nline2: branch-2 content\nline 3: branch-3 content |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-3 | branch-1 |
