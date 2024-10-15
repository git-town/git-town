Feature: shipped parent of a stacked change with conflicting changes

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
    And origin ships the "alpha" branch
    And the current branch is "beta"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git checkout alpha                      |
      | alpha  | git merge --no-edit --ff main           |
      |        | git checkout main                       |
      | main   | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git merge --no-edit --ff origin/beta    |
      |        | git merge --no-edit --ff main           |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress

  Scenario: resolve manually
    When I resolve the conflict in "file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | beta   | git commit --no-edit |
      |        | git push             |
    And the current branch is still "beta"
    And the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, beta |
    And this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |

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
    And the initial branches and lineage exist now
