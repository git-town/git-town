Feature: detaching a branch from a stack with dependent changes

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                         |
      | main   | local, origin | main commit | file      | line 0: main content\nline 1\nline 2\nline 3\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                           |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2\nline 3\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                             |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                               |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3: branch-3 content\nline 4 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-4 | feature | branch-3 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                                                 |
      | branch-4 | local, origin | branch-4 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2: branch-2 content\nline 3: branch-3 content\nline 4: branch-4 content |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-2 | git fetch --prune --tags                                       |
      |          | git checkout branch-3                                          |
      | branch-3 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-1 branch-2 |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
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
      line 4
      """
    And a rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | branch-3 | git rebase --abort    |
      |          | git checkout branch-2 |
    And the initial branches and lineage exist now
    And no rebase is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "file" with:
      """
      line 0: main content
      line 1: branch-1 content
      line 2
      line 3: branch-3 content
      line 4
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | branch-3 | GIT_EDITOR=true git rebase --continue                          |
      |          | git push --force-with-lease                                    |
      |          | git checkout branch-4                                          |
      | branch-4 | git pull                                                       |
      |          | git -c rebase.updateRefs=false rebase --onto branch-3 branch-2 |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    And file "file" now has content:
      """
      line 0: main content
      line 1: branch-1 content
      <<<<<<< HEAD
      line 2
      =======
      line 2: branch-2 content
      >>>>>>> {{ sha-initial-short 'branch-3 commit' }} (branch-3 commit)
      line 3: branch-3 content
      line 4
      """
    And a rebase is now in progress
    When I resolve the conflict in "file" with:
      """
      line 0: main content
      line 1: branch-1 content
      line 2
      line 3: branch-3 content
      line 4
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                    |
      | branch-4 | GIT_EDITOR=true git rebase --continue                      |
      |          | git push --force-with-lease                                |
      |          | git checkout branch-2                                      |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1 |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And a rebase is now in progress
    And file "file" now has content:
      """
      line 0: main content
      <<<<<<< HEAD
      line 1
      line 2
      =======
      line 1: branch-1 content
      line 2: branch-2 content
      >>>>>>> {{ sha-initial-short 'branch-2 commit' }} (branch-2 commit)
      line 3
      line 4
      """
    When I resolve the conflict in "file" with:
      """
      line 0: main content
      line 1
      line 2: branch-2 content
      line 3
      line 4
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT                                                                                               |
      | main     | local, origin | main commit     | file      | line 0: main content\nline 1\nline 2\nline 3\nline 4                                                       |
      | branch-1 | local, origin | branch-1 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2\nline 3\nline 4                                     |
      | branch-3 | local, origin | branch-3 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2\nline 3: branch-3 content\nline 4                   |
      | branch-4 | local, origin | branch-4 commit | file      | line 0: main content\nline 1: branch-1 content\nline 2\nline 3: branch-3 content\nline 4: branch-4 content |
      | branch-2 | local, origin | branch-2 commit | file      | line 0: main content\nline 1\nline 2: branch-2 content\nline 3\nline 4                                     |
    And this lineage exists now
      """
      main
        branch-1
          branch-3
            branch-4
        branch-2
      """
