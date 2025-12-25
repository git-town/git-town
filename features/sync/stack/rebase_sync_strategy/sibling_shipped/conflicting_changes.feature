Feature: conflicting sibling branches, one gets shipped, the other syncs afterwards

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT   |
      | main   | local, origin | main commit | file      | line 1\nline 2 |
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT                     |
      | branch-1 | local, origin | commit 1 | file      | line 1: branch-1 content\nline 2 |
      | branch-2 | local, origin | commit 2 | file      | line 1\nline 2: branch-2 content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                   |
      | branch-2 | git fetch --prune --tags                                                  |
      |          | git checkout main                                                         |
      | main     | git -c rebase.updateRefs=false rebase origin/main                         |
      |          | git checkout branch-2                                                     |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main {{ sha 'main commit' }} |
    And Git Town prints the error:
      """
      CONFLICT (content): Merge conflict in file
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      line 1: branch-1 content
      line 2
      =======
      line 1
      line 2: branch-2 content
      >>>>>>> {{ sha-short 'commit 2' }} (commit 2)
      """
    When I resolve the conflict in "file" with:
      """
      line 1: branch-1 content
      line 2: branch-2 content
      """
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, branch-1, branch-2 |
      | origin     | main, branch-2           |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                                       |
      | main     | local, origin | main commit | file      | line 1\nline 2                                     |
      |          |               | commit 1    | file      | line 1: branch-1 content\nline 2                   |
      | branch-1 | local         | commit 1    | file      | line 1: branch-1 content\nline 2                   |
      | branch-2 | local, origin | commit 2    | file      | line 1: branch-1 content\nline 2: branch-2 content |
    And this lineage exists now
      """
      main
        branch-1
        branch-2
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-2 | git rebase --abort                       |
      |          | git checkout main                        |
      | main     | git reset --hard {{ sha 'main commit' }} |
      |          | git checkout branch-2                    |
    And the initial branches and lineage exist now
