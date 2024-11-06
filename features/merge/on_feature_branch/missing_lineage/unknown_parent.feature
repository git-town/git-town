Feature: merging with missing lineage

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE   | PARENT | LOCATIONS     |
      | alpha | (none) |        | local, origin |
      | beta  | (none) |        | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And the current branch is "beta"
    When I run "git-town merge" and enter into the dialog:
      | DIALOG                          | KEYS       |
      | select parent branch for "beta" | down enter |
      | select parent branch for "beta" | enter      |

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git push                              |
      |        | git branch -D alpha                   |
      |        | git push origin :alpha                |
    And the current branch is still "beta"
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
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git push -u origin alpha                             |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
