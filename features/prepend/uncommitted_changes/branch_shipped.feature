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
    And an uncommitted file
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                           |
      | child  | git add -A                        |
      |        | git stash                         |
      |        | git checkout main                 |
      | main   | git rebase origin/main            |
      |        | git checkout parent               |
      | parent | git merge --no-edit origin/parent |
      |        | git merge --no-edit main          |
      |        | git checkout child                |
      | child  | git merge --no-edit origin/child  |
      |        | git merge --no-edit parent        |
      |        | git push                          |
      |        | git checkout -b new parent        |
      | new    | git stash pop                     |
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, child, parent      |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | new    |
      | new    | parent |
      | parent | main   |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout child                              |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git branch -D new                               |
      |        | git stash pop                                   |
    And the current branch is now "child"
    And the branches are now
      | REPOSITORY    | BRANCHES            |
      | local, origin | main, child, parent |
    And the initial lineage exists
    And the uncommitted file still exists
