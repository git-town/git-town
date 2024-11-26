Feature: delete the current parked branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | parked  | parked  | main   | local, origin |
      | feature | feature | main   | local, origin |
    And the current branch is "parked"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | parked  | local, origin | parked commit  |
    And an uncommitted file
    And the current branch is "parked" and the previous branch is "feature"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                          |
      | parked  | git fetch --prune --tags                         |
      |         | git push origin :parked                          |
      |         | git add -A                                       |
      |         | git commit -m "Committing WIP on deleted branch" |
      |         | git checkout feature                             |
      | feature | git branch -D parked                             |
    And the current branch is now "feature"
    And no uncommitted files exist now
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                        |
      | feature | git push origin {{ sha 'parked commit' }}:refs/heads/parked    |
      |         | git branch parked {{ sha 'Committing WIP on deleted branch' }} |
      |         | git checkout parked                                            |
      | parked  | git reset --soft HEAD~1                                        |
    And the current branch is now "parked"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
    And branch "parked" is now parked
