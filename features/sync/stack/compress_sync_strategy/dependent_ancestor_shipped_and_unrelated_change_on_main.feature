Feature: shipped the head branch of a synced stack with dependent changes while main also received independent updates

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | file      | beta content |
    And Git Town setting "sync-feature-strategy" is "compress"
    And origin ships the "alpha" branch
    And I add commit "additional commit" to the "main" branch
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git push                                |
      |        | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git merge --no-edit --ff main           |
      |        | git checkout --ours file                |
      |        | git add file                            |
      |        | git commit --no-edit                    |
      |        | git merge --no-edit --ff origin/beta    |
      |        | git reset --soft main                   |
      |        | git commit -m "beta commit"             |
      |        | git push --force-with-lease             |
    And the current branch is still "beta"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT  |
      | main   | local, origin | alpha commit      | file      | alpha content |
      |        |               | additional commit | new_file  |               |
      | beta   | local, origin | beta commit       | file      | beta content  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE           | FILE NAME | FILE CONTENT  |
      | main   | local, origin | alpha commit      | file      | alpha content |
      |        |               | additional commit | new_file  |               |
      | alpha  | local         | alpha commit      | file      | alpha content |
      | beta   | local, origin | beta commit       | file      | beta content  |
      |        | origin        | alpha commit      | file      | alpha content |
    And the initial branches and lineage exist now
