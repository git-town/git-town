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
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                 |
      | beta   | git fetch --prune --tags                                |
      |        | git checkout main                                       |
      | main   | git -c rebase.updateRefs=false rebase origin/main       |
      |        | git checkout beta                                       |
      | beta   | git pull                                                |
      |        | git -c rebase.updateRefs=false rebase --onto main alpha |
      |        | git push --force-with-lease                             |
      |        | git branch -D alpha                                     |
    And no rebase is now in progress
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME      | FILE CONTENT |
      | main   | local, origin | alpha commit 1 | favorite-fruit | peach        |
      | beta   | local, origin | beta commit 1  | favorite-pizza | pepperoni    |
      |        |               | beta commit 2  | favorite-pizza | pineapple    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                             |
      | beta   | git reset --hard {{ sha 'beta commit 2' }}          |
      |        | git push --force-with-lease --force-if-includes     |
      |        | git checkout main                                   |
      | main   | git reset --hard {{ sha 'initial commit' }}         |
      |        | git branch alpha {{ sha-initial 'alpha commit 2' }} |
      |        | git checkout beta                                   |
    And the initial branches and lineage exist now
    And the initial commits exist now
