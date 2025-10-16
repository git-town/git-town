Feature: prepending to a branch whose parent was shipped and the local branch deleted manually

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
    And I ran "git branch -d parent"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | child  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout child                                |
      | child  | git merge --no-edit --ff main                     |
      |        | git push                                          |
      |        | git checkout -b new main                          |
    And this lineage exists now
      """
      main
        new
          child
      """
    And the branches are now
      | REPOSITORY | BRANCHES         |
      | local      | main, child, new |
      | origin     | main, child      |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout child                              |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git checkout child                              |
      | child  | git branch -D new                               |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
