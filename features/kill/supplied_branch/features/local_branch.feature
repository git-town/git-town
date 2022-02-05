Feature: local branch

  Background:
    Given my repo does not have a remote origin
    And my repo has the local feature branches "dead" and "other"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE      |
      | dead   | local    | dead commit  |
      | other  | local    | other commit |
    And I am on the "dead" branch
    And my workspace has an uncommitted file
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | dead   | git add -A                  |
      |        | git commit -m "WIP on dead" |
      |        | git checkout main           |
      | main   | git branch -D dead          |
    And I am now on the "main" branch
    And my repo doesn't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
    And my repo now has the commits
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
    And I am now on the "dead" branch
    And my workspace has the uncommitted file again
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
