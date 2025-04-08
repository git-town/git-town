Feature: syncing a branch whose parent with independent changes was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "parent" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | child  | git fetch --prune --tags                       |
      |        | git checkout main                              |
      | main   | git rebase origin/main --no-update-refs        |
      |        | git checkout child                             |
      | child  | git pull                                       |
      |        | git rebase --onto main parent --no-update-refs |
      |        | git push --force-with-lease                    |
      |        | git branch -D parent                           |
    And Git Town prints:
      """
      deleted branch "parent"
      """
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                |
      | child  | git reset --hard {{ sha 'child commit' }}              |
      |        | git push --force-with-lease --force-if-includes        |
      |        | git checkout main                                      |
      | main   | git reset --hard {{ sha 'initial commit' }}            |
      |        | git branch parent {{ sha-before-run 'parent commit' }} |
      |        | git checkout child                                     |
    And the initial branches and lineage exist now
