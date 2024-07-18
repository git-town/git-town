@smoke
Feature: delete another than the current branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | good | feature | main   | local, origin |
      | dead | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git push origin :dead    |
      |        | git branch -D dead       |
      |        | git stash pop            |
    And the current branch is still "good"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, good |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | conflicting commit |
      | good   | local, origin | good commit        |
    And this lineage exists now
      | BRANCH | PARENT |
      | good   | main   |

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
