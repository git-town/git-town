Feature: shipped the head branch of a synced stack with dependent changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
    # | alpha  | local, origin | alpha commit 2 | file_2      | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | file      | beta content |
    # | beta   | local, origin | beta commit 2 | file_2      | beta content |
    And the current branch is "beta"
    And origin ships the "alpha" branch using the "squash-merge" ship-strategy
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git rebase --onto main alpha            |
      |        | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git merge --no-edit --ff main           |
      |        | git checkout --ours file                |
      |        | git add file                            |
      |        | git commit --no-edit                    |
      |        | git merge --no-edit --ff origin/beta    |
      |        | git push                                |
    And the current branch is still "beta"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                       | FILE NAME | FILE CONTENT  |
      | main   | local, origin | alpha commit                  | file      | alpha content |
      | beta   | local, origin | alpha commit                  | file      | alpha content |
      |        |               | beta commit                   | file      | beta content  |
      |        |               | Merge branch 'main' into beta |           |               |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git checkout main                                    |
      | main   | git reset --hard {{ sha 'initial commit' }}          |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git checkout beta                                    |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | main   | origin        | alpha commit | file      | alpha content |
      | alpha  | local         | alpha commit | file      | alpha content |
      | beta   | local, origin | beta commit  | file      | beta content  |
      |        | origin        | alpha commit | file      | alpha content |
    And the initial branches and lineage exist now
