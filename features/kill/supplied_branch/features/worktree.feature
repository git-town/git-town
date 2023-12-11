Feature: delete a branch that is checked out in another worktree

  Background:
    Given the feature branches "good" and "dead"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And branch "dead" is checked out in another worktree
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill dead"

  Scenario: result
    Then it prints the error:
      """
      I cannot kill this branch because it is checked out in another workspace
      """
    And it runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | good   | git add -A                                  |
      |        | git stash                                   |
      |        | git branch dead {{ sha 'dead-end commit' }} |
      |        | git push -u origin dead                     |
      |        | git stash pop                               |
    And the current branch is still "good"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
