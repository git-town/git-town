Feature: multiple shipped parent branches in a lineage

  Background:
    Given a feature branch "feature-1"
    And a feature branch "feature-2" as a child of "feature-1"
    And a feature branch "feature-3" as a child of "feature-2"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
      | feature-2 | local, origin | feature-2 commit | feature-2-file | feature 2 content |
      | feature-3 | local, origin | feature-3 commit | feature-3-file | feature 3 content |
    And origin ships the "feature-1" branch
    And origin ships the "feature-2" branch
    And the current branch is "feature-3"
    When I run "git-town sync"

  @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                              |
      | feature-3 | git fetch --prune --tags             |
      |           | git checkout feature-1               |
      | feature-1 | git merge --no-edit main             |
      |           | git checkout main                    |
      | main      | git branch -d feature-1              |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit main             |
      |           | git checkout main                    |
      | main      | git branch -d feature-2              |
      |           | git rebase origin/main               |
      |           | git checkout feature-3               |
      | feature-3 | git merge --no-edit origin/feature-3 |
      |           | git merge --no-edit main             |
    And it prints:
      """
      deleted branch "feature-1"
      """
    And it prints:
      """
      deleted branch "feature-2"
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
