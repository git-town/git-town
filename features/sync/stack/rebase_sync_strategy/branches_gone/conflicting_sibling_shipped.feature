Feature: conflicting sibling branches, one gets shipped, the other syncs afterwards
# TODO: this should remove the local branch-1, but doesn't

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1 | file      | content 1    |
      | branch-2 | local, origin | commit 2 | file      | content 2    |
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | branch-2 | git fetch --prune --tags                          |
      |          | git checkout main                                 |
      | main     | git -c rebase.updateRefs=false rebase origin/main |
      |          | git checkout branch-2                             |
      | branch-2 | git -c rebase.updateRefs=false rebase main        |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And file "file" now has content:
      """
      <<<<<<< HEAD
      content 1
      =======
      content 2
      >>>>>>> {{ sha-short 'commit 2' }} (commit 2)
      """
    When I resolve the conflict in "file" with "content 1 and 2"
    And I run "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, branch-1, branch-2 |
      | origin     | main, branch-2           |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT    |
      | main     | local, origin | commit 1 | file      | content 1       |
      | branch-1 | local         | commit 1 | file      | content 1       |
      | branch-2 | local, origin | commit 2 | file      | content 1 and 2 |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-1 | main   |
      | branch-2 | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                          |
      | branch-2 | git reset --hard {{ sha-initial 'commit 4' }}    |
      |          | git push --force-with-lease --force-if-includes  |
      |          | git checkout main                                |
      | main     | git reset --hard {{ sha 'initial commit' }}      |
      |          | git branch branch-1 {{ sha-initial 'commit 2' }} |
      |          | git checkout branch-2                            |
    And the initial branches and lineage exist now
