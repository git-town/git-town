Feature: using the "compress" strategy, sync a branch whose parent was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE        |
      | parent | local, origin | parent commit  |
      | child  | local, origin | child commit 1 |
      | child  | local, origin | child commit 2 |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And origin ships the "parent" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | child  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git branch -D parent                              |
      |        | git checkout child                                |
      | child  | git merge --no-edit --ff main                     |
      |        | git reset --soft main --                          |
      |        | git commit -m "child commit 1"                    |
      |        | git push --force-with-lease                       |
    And Git Town prints:
      """
      deleted branch parent
      """
    And this lineage exists now
      """
      main
        child
      """
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | parent commit  |
      | child  | local, origin | child commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | child  | git reset --hard {{ sha-initial 'child commit 2' }} |
      |        | git push --force-with-lease --force-if-includes     |
      |        | git checkout main                                   |
      | main   | git reset --hard {{ sha 'initial commit' }}         |
      |        | git branch parent {{ sha-initial 'parent commit' }} |
      |        | git checkout child                                  |
    And the initial branches and lineage exist now
