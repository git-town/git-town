Feature: prepend a branch to a branch that was shipped at the remote

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
    And origin ships the "child" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | child  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout parent                               |
      | parent | git merge --no-edit --ff main                     |
      |        | git push                                          |
      |        | git checkout main                                 |
      | main   | git branch -D child                               |
      |        | git checkout -b new parent                        |
    And Git Town prints:
      """
      deleted branch "child"
      """
    And Git Town prints:
      """
      branch "new" is now a child of "parent"
      """
    And this lineage exists now
      """
      main
        parent
          new
      """
    And the branches are now
      | REPOSITORY | BRANCHES          |
      | local      | main, new, parent |
      | origin     | main, parent      |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout parent                             |
      | parent | git reset --hard {{ sha 'parent commit' }}      |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch child {{ sha 'child commit' }}       |
      |        | git checkout child                              |
      | child  | git branch -D new                               |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES            |
      | local      | main, child, parent |
      | origin     | main, parent        |
