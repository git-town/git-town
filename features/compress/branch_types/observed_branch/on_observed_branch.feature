Feature: does not compress an active observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE    | FILE NAME  | FILE CONTENT |
      | observed | local, origin | observed 1 | observed_1 | observed 1   |
      |          |               | observed 2 | observed_2 | observed 2   |
    And the branches
      | NAME  | TYPE    | PARENT   | LOCATIONS     |
      | child | feature | observed | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "observed"
    When I run "git-town compress --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | observed | git fetch --prune --tags                        |
      |          | git checkout child                              |
      | child    | git reset --soft observed --                    |
      |          | git commit -m "child 1"                         |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout observed                           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE    |
      | observed | local, origin | observed 1 |
      |          |               | observed 2 |
      | child    | local, origin | child 1    |
    And file "observed_1" still has content "observed 1"
    And file "observed_2" still has content "observed 2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | observed | git checkout child                              |
      | child    | git reset --hard {{ sha 'child 2' }}            |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout observed                           |
    And the initial branches and lineage exist now
    And the initial commits exist now
