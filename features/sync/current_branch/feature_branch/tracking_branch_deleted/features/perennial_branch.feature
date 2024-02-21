Feature: sync perennial branch that was deleted at the remote

  Background:
    Given the perennial branches "feature-1" and "feature-2"
    And a feature branch "feature-1a" as a child of "feature-1"
    And a feature branch "feature-1b" as a child of "feature-1"
    And a feature branch "feature-2a" as a child of "feature-2"
    And origin deletes the "feature-1" branch
    And the current branch is "feature-1"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | feature-1 | git fetch --prune --tags |
      |           | git checkout main        |
      | main      | git branch -D feature-1  |
      |           | git push --tags          |
    And it prints:
      """
      deleted branch "feature-1"
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES                                            |
      | local, origin | main, feature-1a, feature-1b, feature-2, feature-2a |
    And the perennial branches are now "feature-2"
    And this branch lineage exists now
      | BRANCH     | PARENT    |
      | feature-1a | main      |
      | feature-1b | main      |
      | feature-2a | feature-2 |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git branch feature-1 {{ sha 'initial commit' }} |
      |        | git checkout feature-1                          |
    And the current branch is now "feature-1"
    And the initial branches and lineage exist
    And the perennial branches are now "feature-1" and "feature-2"
