Feature: delete a local branch

  Background:
    Given the current branch is a local feature branch "local"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | local  | local    | local commit |
    And an uncommitted file
    When I run "git-town kill"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                      |
      | local  | git fetch --prune --tags     |
      |        | git add -A                   |
      |        | git commit -m "WIP on local" |
      |        | git checkout main            |
      | main   | git branch -D local          |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no branch hierarchy exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git branch local {{ sha 'WIP on local' }} |
      |        | git checkout local                        |
      | local  | git reset {{ sha 'local commit' }}        |
    And the current branch is now "local"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branches and hierarchy exist
