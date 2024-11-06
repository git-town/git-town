Feature: merging a branch with a parent that has conflicting changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME        | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | conflicting_file | alpha content |
      | beta   | local, origin | beta commit  | conflicting_file | beta content  |
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | beta   | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff alpha        |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH | COMMAND                              |
      | beta   | git commit --no-edit                 |
      |        | git merge --no-edit --ff origin/beta |
      |        | git push                             |
      |        | git branch -D alpha                  |
      |        | git push origin :alpha               |
    And the current branch is still "beta"
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | beta   | local, origin | beta commit                    |
      |        |               | alpha commit                   |
      |        |               | Merge branch 'alpha' into beta |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | beta   | git merge --abort |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
