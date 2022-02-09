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
    And my repo doesn't have any uncommitted files
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE      |
      | other  | local    | other commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git branch dead {{ sha 'WIP on dead' }} |
      |        | git checkout dead                       |
      | dead   | git reset {{ sha 'dead commit' }}       |
    And the current branch is now "dead"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branches and hierarchy exist
