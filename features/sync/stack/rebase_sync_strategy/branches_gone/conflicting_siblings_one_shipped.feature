Feature: conflicting sibling branches, some shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
      | branch-2 | feature | main   | local, origin |
      | branch-3 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | branch-1 | local, origin | commit 1  | file_a    | content 1    |
      | branch-2 | local, origin | commit 2  | file_b    | content 2    |
      | branch-3 | local, origin | commit 3a | file_a    | content 3a   |
      |          | local, origin | commit 3b | file_b    | content 3b   |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And origin ships the "branch-2" branch using the "squash-merge" ship-strategy
    And the current branch is "main"
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                           |
      | main     | git fetch --prune --tags                          |
      |          | git -c rebase.updateRefs=false rebase origin/main |
      |          | git branch -D branch-1                            |
      |          | git branch -D branch-2                            |
      |          | git checkout branch-3                             |
      | branch-3 | git -c rebase.updateRefs=false rebase main        |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file_a
      """
    And file "file_a" now has content:
      """
      <<<<<<< HEAD
      content 1
      =======
      content 3a
      >>>>>>> {{ sha-short 'commit 3a' }} (commit 3a)
      """
    When I resolve the conflict in "file_a" with "content 3a"
    And I run "git town continue"
    Then Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file_b
      """
    And file "file_b" now has content:
      """
      <<<<<<< HEAD
      content 2
      =======
      content 3b
      >>>>>>> {{ sha-short 'commit 3b' }} (commit 3b)
      """
    When I resolve the conflict in "file_b" with "content 3b"
    And I run "git town continue"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-3 | GIT_EDITOR=true git rebase --continue           |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout main                               |
      | main     | git push --tags                                 |
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, branch-3 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | main     | local, origin | commit 1  | file_a    | content 1    |
      |          |               | commit 2  | file_b    | content 2    |
      | branch-3 | local, origin | commit 3a | file_a    | content 3a   |
      |          |               | commit 3b | file_b    | content 3b   |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-3 | main   |

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
