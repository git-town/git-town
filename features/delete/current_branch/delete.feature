@smoke
Feature: delete the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | current | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the current branch is "current"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | current | local, origin | current commit |
      | other   | local, origin | other commit   |
    And an uncommitted file
    And the current branch is "current" and the previous branch is "other"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                          |
      | current | git fetch --prune --tags                         |
      |         | git push origin :current                         |
      |         | git add -A                                       |
      |         | git commit -m "Committing WIP for git town undo" |
      |         | git checkout other                               |
      | other   | git rebase --onto main current                   |
      |         | git branch -D current                            |
    And the current branch is now "other"
    And no uncommitted files exist now
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | other  | local, origin | other commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                         |
      | other   | git push origin {{ sha 'current commit' }}:refs/heads/current   |
      |         | git branch current {{ sha 'Committing WIP for git town undo' }} |
      |         | git checkout current                                            |
      | current | git reset --soft HEAD~1                                         |
    And the current branch is now "current"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now
