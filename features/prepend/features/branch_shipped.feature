Feature: prepend a branch to a branch that was shipped at the remote

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "child" branch
    And the current branch is "child"
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                           |
      | child  | git fetch --prune --tags          |
      |        | git checkout main                 |
      | main   | git rebase origin/main            |
      |        | git checkout parent               |
      | parent | git merge --no-edit origin/parent |
      |        | git merge --no-edit main          |
      |        | git push                          |
      |        | git checkout child                |
      | child  | git merge --no-edit parent        |
      |        | git checkout parent               |
      | parent | git branch -D child               |
      |        | git branch new parent             |
      |        | git checkout new                  |
    And it prints:
      """
      deleted branch "child"
      """
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES          |
      | local      | main, new, parent |
      | origin     | main, parent      |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | new    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout parent                             |
      | parent | git reset --hard {{ sha 'parent commit' }}      |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch child {{ sha 'child commit' }}       |
      |        | git checkout child                              |
      | child  | git branch -D new                               |
    And the current branch is now "child"
    And the initial branches and lineage exist
