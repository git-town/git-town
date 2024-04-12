Feature: does not compress contribution branches in the stack

  Background:
    Given contribution branch "contribution" with these commits
      | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT   |
      | local, origin | contribution 1 | contribution_1 | contribution 1 |
      |               | contribution 2 | contribution_2 | contribution 2 |
      |               | contribution 3 | contribution_3 | contribution 3 |
    And feature branch "child" as a child of "contribution" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | child 1 | child_1   | child 1      |
      |               | child 2 | child_2   | child 2      |
      |               | child 3 | child_3   | child 3      |
    And the current branch is "child"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --soft contribution                   |
      |        | git commit -m "child 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "child"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE        |
      | child        | local, origin | contribution 1 |
      |              |               | contribution 2 |
      |              |               | contribution 3 |
      |              |               | child 1        |
      | contribution | local, origin | contribution 1 |
      |              |               | contribution 2 |
      |              |               | contribution 3 |
    And file "contribution_1" still has content "contribution 1"
    And file "contribution_2" still has content "contribution 2"
    And file "contribution_3" still has content "contribution 3"
    And file "child_1" still has content "child 1"
    And file "child_2" still has content "child 2"
    And file "child_3" still has content "child 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --hard {{ sha 'child 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
