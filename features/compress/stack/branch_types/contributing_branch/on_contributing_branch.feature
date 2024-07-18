Feature: does not compress the active contribution branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT   |
      | contribution | local, origin | contribution 1 | contribution_1 | contribution 1 |
      |              |               | contribution 2 | contribution_2 | contribution 2 |
    And the branches
      | NAME  | TYPE    | PARENT       | LOCATIONS     |
      | child | feature | contribution | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "contribution"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                                         |
      | contribution | git fetch --prune --tags                        |
      |              | git add -A                                      |
      |              | git stash                                       |
      |              | git checkout child                              |
      | child        | git reset --soft contribution                   |
      |              | git commit -m "child 1"                         |
      |              | git push --force-with-lease --force-if-includes |
      |              | git checkout contribution                       |
      | contribution | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "contribution"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE        |
      | child        | local, origin | contribution 1 |
      |              |               | contribution 2 |
      |              |               | child 1        |
      | contribution | local, origin | contribution 1 |
      |              |               | contribution 2 |
    And file "contribution_1" still has content "contribution 1"
    And file "contribution_2" still has content "contribution 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                         |
      | contribution | git add -A                                      |
      |              | git stash                                       |
      |              | git checkout child                              |
      | child        | git reset --hard {{ sha 'child 2' }}            |
      |              | git push --force-with-lease --force-if-includes |
      |              | git checkout contribution                       |
      | contribution | git stash pop                                   |
    And the current branch is still "contribution"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
