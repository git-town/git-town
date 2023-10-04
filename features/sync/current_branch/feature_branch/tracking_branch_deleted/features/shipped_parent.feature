Feature: a parent branch of the current feature branch was shipped

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin deletes the "parent" branch
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | child  | git fetch --prune --tags         |
      |        | git checkout main                |
      | main   | git rebase origin/main           |
      |        | git branch -D parent             |
      |        | git checkout child               |
      | child  | git merge --no-edit origin/child |
      |        | git merge --no-edit main         |
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | child  | git branch parent {{ sha 'parent commit' }} |
    And the current branch is still "child"
    And the initial branches and hierarchy exist
