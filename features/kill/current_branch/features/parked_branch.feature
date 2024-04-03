Feature: delete the current parked branch

  Background:
    Given the current branch is a parked branch "parked"
    And a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
      | parked  | local, origin | parked commit  |
    And an uncommitted file
    And the current branch is "parked" and the previous branch is "feature"
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | parked  | git fetch --prune --tags      |
      |         | git push origin :parked       |
      |         | git add -A                    |
      |         | git commit -m "WIP on parked" |
      |         | git checkout feature          |
      | feature | git branch -D parked          |
    And the current branch is now "feature"
    And no uncommitted files exist
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
    Then it runs the commands
      | BRANCH  | COMMAND                                                     |
      | feature | git push origin {{ sha 'parked commit' }}:refs/heads/parked |
      |         | git branch parked {{ sha 'WIP on parked' }}                 |
      |         | git checkout parked                                         |
      | parked  | git reset --soft HEAD~1                                     |
    And the current branch is now "parked"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
    And branch "parked" is now parked
