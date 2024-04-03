Feature: local branch

  Background:
    Given my repo does not have an origin
    And the local feature branches "dead" and "other"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | dead   | local    | dead commit  |
      | other  | local    | other commit |
    And the current branch is "dead"
    And an uncommitted file
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | dead   | git add -A                  |
      |        | git commit -m "WIP on dead" |
      |        | git checkout main           |
      | main   | git branch -D dead          |
    And the current branch is now "main"
    And no uncommitted files exist
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git branch dead {{ sha 'WIP on dead' }} |
      |        | git checkout dead                       |
      | dead   | git reset --soft HEAD~1                 |
    And the current branch is now "dead"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
