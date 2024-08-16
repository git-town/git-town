Feature: syncing a branch whose parent was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "parent" branch
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main                          |
      |        | git checkout parent                             |
      | parent | git rebase main                                 |
      |        | git checkout main                               |
      | main   | git branch -D parent                            |
      |        | git checkout child                              |
      | child  | git rebase main                                 |
      |        | git push --force-with-lease --force-if-includes |
    And it prints:
      """
      deleted branch "parent"
      """
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                |
      | child  | git reset --hard {{ sha 'child commit' }}              |
      |        | git push --force-with-lease --force-if-includes        |
      |        | git checkout main                                      |
      | main   | git reset --hard {{ sha 'initial commit' }}            |
      |        | git branch parent {{ sha-before-run 'parent commit' }} |
      |        | git checkout child                                     |
    And the current branch is still "child"
    And the initial branches and lineage exist
