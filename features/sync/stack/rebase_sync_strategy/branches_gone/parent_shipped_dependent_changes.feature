Feature: syncing a branch whose parent with dependent changes was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT  |
      | branch-1 | local, origin | commit 1 | file.txt  | new content   |
      |          | local, origin | commit 2 | file.txt  | new content 2 |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT  |
      | branch-2 | local, origin | commit 3 | file.txt  | new content 3 |
      |          | local, origin | commit 4 | file.txt  | new content 4 |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "branch-1" branch using the "squash-merge" ship-strategy
    And the current branch is "branch-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                          |
      | branch-2 | git fetch --prune --tags                         |
      |          | git checkout main                                |
      | main     | git rebase origin/main --no-update-refs          |
      |          | git checkout branch-2                            |
      | branch-2 | git pull                                         |
      |          | git rebase --onto main branch-1 --no-update-refs |
      |          | git push --force-with-lease                      |
      |          | git branch -D branch-1                           |
    And Git Town prints:
      """
      deleted branch "branch-1"
      """
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, branch-2 |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT  |
      | main     | local, origin | commit 1 | file.txt  | new content 2 |
      | branch-2 | local, origin | commit 3 | file.txt  | new content 3 |
      |          |               | commit 4 | file.txt  | new content 4 |
    And this lineage exists now
      | BRANCH   | PARENT |
      | branch-2 | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                             |
      | branch-2 | git reset --hard {{ sha-before-run 'commit 4' }}    |
      |          | git push --force-with-lease --force-if-includes     |
      |          | git checkout main                                   |
      | main     | git reset --hard {{ sha 'initial commit' }}         |
      |          | git branch branch-1 {{ sha-before-run 'commit 2' }} |
      |          | git checkout branch-2                               |
    And the initial branches and lineage exist now
