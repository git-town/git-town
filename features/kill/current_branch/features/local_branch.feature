Feature: delete a local branch

  Background:
    And a local feature branch "local"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | local  | local    | local commit |
    And I am on the "local" branch
    And my workspace has an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | local  | git fetch --prune --tags     |
      |        | git add -A                   |
      |        | git commit -m "WIP on local" |
      |        | git checkout main            |
      | main   | git branch -D local          |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And Git Town is now aware of no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git branch local {{ sha 'WIP on local' }} |
      |        | git checkout local                        |
      | local  | git reset {{ sha 'local commit' }}        |
    And I am now on the "local" branch
    And my workspace still contains my uncommitted file
    And now the initial commits exist
    And my repo now has its initial branches and branch hierarchy
