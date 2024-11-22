Feature: stack that changes the same file in multiple commits per branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT |
      | alpha  | local, origin | alpha commit 1 | favorite-fruit | apple        |
      | alpha  | local, origin | alpha commit 2 | favorite-fruit | peach        |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME      | FILE CONTENT |
      | beta   | local, origin | beta commit 1 | favorite-pizza | pepperoni    |
      | beta   | local, origin | beta commit 2 | favorite-pizza | pineapple    |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "beta"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git rebase --onto main alpha            |
      |        | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git rebase main --no-update-refs        |
    And no rebase is in progress
    And the current branch is still "beta"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT   |
      | main   | local, origin | alpha commit 1 | favorite-fruit | peach          |
      | beta   | local, origin | alpha commit 1 | favorite-fruit | resolved apple |
      |        |               | alpha commit 2 | favorite-fruit | resolved peach |
      |        |               | beta commit 1  | favorite-pizza | pepperoni      |
      |        |               | beta commit 2  | favorite-pizza | pineapple      |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                |
      | beta   | git rebase --abort                                     |
      |        | git checkout main                                      |
      | main   | git reset --hard {{ sha 'initial commit' }}            |
      |        | git branch alpha {{ sha-before-run 'alpha commit 2' }} |
      |        | git checkout beta                                      |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT |
      | main   | origin        | alpha commit 1 | favorite-fruit | peach        |
      | alpha  | local         | alpha commit 1 | favorite-fruit | apple        |
      |        |               | alpha commit 2 | favorite-fruit | peach        |
      | beta   | local, origin | beta commit 1  | favorite-pizza | pepperoni    |
      |        |               | beta commit 2  | favorite-pizza | pineapple    |
      |        | origin        | alpha commit 1 | favorite-fruit | apple        |
      |        |               | alpha commit 2 | favorite-fruit | peach        |
    And the initial branches and lineage exist now
