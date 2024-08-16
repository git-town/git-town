Feature: append a branch to a branch whose parent was shipped on the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent |               |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "parent" branch
    And the current branch is "child"
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | child  | git fetch --prune --tags              |
      |        | git checkout main                     |
      | main   | git rebase origin/main                |
      |        | git checkout parent                   |
      | parent | git merge --no-edit --ff main         |
      |        | git checkout main                     |
      | main   | git branch -D parent                  |
      |        | git checkout child                    |
      | child  | git merge --no-edit --ff origin/child |
      |        | git merge --no-edit --ff main         |
      |        | git push                              |
      |        | git checkout -b new                   |
    And it prints:
      """
      deleted branch "parent"
      """
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES         |
      | local      | main, child, new |
      | origin     | main, child      |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | new    | child  |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | main   | local, origin | parent commit                  |
      | child  | local, origin | child commit                   |
      |        |               | parent commit                  |
      |        |               | Merge branch 'main' into child |
      | new    | local         | child commit                   |
      |        |               | parent commit                  |
      |        |               | Merge branch 'main' into child |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                |
      | new    | git checkout child                                     |
      | child  | git reset --hard {{ sha 'child commit' }}              |
      |        | git push --force-with-lease --force-if-includes        |
      |        | git checkout main                                      |
      | main   | git reset --hard {{ sha 'initial commit' }}            |
      |        | git branch parent {{ sha-before-run 'parent commit' }} |
      |        | git checkout child                                     |
      | child  | git branch -D new                                      |
    And the current branch is still "child"
    And the initial branches and lineage exist
