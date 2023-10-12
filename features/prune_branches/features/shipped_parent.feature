Feature: a parent branch of a local branch was shipped

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | child  | local, origin | child commit |
    And origin ships the "parent" branch
    And the current branch is "child"
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | child  | git fetch --prune --tags |
      |        | git checkout main        |
      | main   | git rebase origin/main   |
      |        | git checkout parent      |
      | parent | git merge --no-edit main |
      |        | git checkout main        |
      | main   | git branch -d parent     |
      |        | git checkout child       |
    And it prints:
      """
      deleted branch "parent"
      """
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
      | BRANCH | COMMAND                                      |
      | child  | git branch parent {{ sha 'Initial commit' }} |
    And the current branch is still "child"
    And the initial branches and hierarchy exist
