Feature: sync perennial branch that was deleted at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT    | LOCATIONS     |
      | perennial | perennial |           | local, origin |
      | feature-1 | feature   | perennial | local, origin |
      | feature-2 | feature   | perennial | local, origin |
    And origin deletes the "perennial" branch
    And the current branch is "perennial" and the previous branch is "main"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
      |           | git checkout feature-1   |
      | feature-1 | git branch -D perennial  |
      |           | git checkout main        |
      | main      | git push --tags          |
    And Git Town prints:
      """
      deleted branch perennial
      """
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-1, feature-2 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git branch perennial {{ sha 'initial commit' }} |
      |        | git checkout perennial                          |
    And the initial branches and lineage exist now
    And branch "perennial" now has type "perennial"
