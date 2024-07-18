Feature: sync perennial branch that was deleted at the remote

  Background:
    Given a Git repo clone
    And the branches
      | NAME       | TYPE      | PARENT    | LOCATIONS     |
      | feature-1  | perennial |           | local, origin |
      | feature-2  | perennial |           | local, origin |
      | feature-1a | feature   | feature-1 | local, origin |
      | feature-1b | feature   | feature-1 | local, origin |
      | feature-2a | feature   | feature-2 | local, origin |
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
    And this lineage exists now
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
