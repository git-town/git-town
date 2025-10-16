Feature: delete a branch that has an overridden branch type

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
      | other        | feature      | main   | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
      | other        | local, origin | other commit        |
    And Git setting "git-town-branch.contribution.branchtype" is "feature"
    And the current branch is "contribution"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                       |
      | contribution | git fetch --prune --tags      |
      |              | git push origin :contribution |
      |              | git checkout other            |
      | other        | git branch -D contribution    |
    And Git setting "git-town-branch.contribution.branchtype" now doesn't exist
    And this lineage exists now
      """
      main
        other
      """
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                 |
      | other  | git branch contribution {{ sha 'contribution commit' }} |
      |        | git push -u origin contribution                         |
      |        | git checkout contribution                               |
    And Git setting "git-town-branch.contribution.branchtype" is now "feature"
    And the initial branches and lineage exist now
    And the initial commits exist now
