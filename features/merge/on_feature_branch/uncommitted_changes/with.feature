Feature: merging a branch with uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And the current branch is "beta"
    And an uncommitted file
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git add -A                            |
      |        | git stash                             |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git push                              |
      |        | git branch -D alpha                   |
      |        | git push origin :alpha                |
      |        | git stash pop                         |
    And the current branch is still "beta"
    And the uncommitted file still exists
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        | FILE NAME  | FILE CONTENT  |
      | beta   | local, origin | beta commit                    | beta-file  | beta content  |
      |        |               | alpha commit                   | alpha-file | alpha content |
      |        |               | Merge branch 'alpha' into beta |            |               |
    And these committed files exist now
      | BRANCH | NAME       | CONTENT       |
      | beta   | alpha-file | alpha content |
      |        | beta-file  | beta content  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git add -A                                           |
      |        | git stash                                            |
      |        | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git push -u origin alpha                             |
      |        | git stash pop                                        |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
