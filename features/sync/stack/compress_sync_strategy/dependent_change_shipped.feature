Feature: shipped the head branch of a synced stack with dependent changes

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
    And the current branch is "beta"
    And Git Town setting "sync-feature-strategy" is "compress"
    And origin ships the "alpha" branch
    When I run "git-town sync"

  # TODO: automatically resolve this phantom merge conflict
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git merge --no-edit --ff origin/beta    |
      |        | git merge --no-edit --ff main           |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | beta   | git commit --no-edit        |
      |        | git reset --soft main       |
      |        | git commit -m "beta commit" |
      |        | git push --force-with-lease |
    And the current branch is still "beta"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT     |
      | main   | local, origin | alpha commit | file      | alpha content    |
      | beta   | local, origin | beta commit  | file      | resolved content |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git merge --abort                                    |
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