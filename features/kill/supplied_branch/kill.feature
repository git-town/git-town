Feature: delete another than the current branch

  Background:
    Given the feature branches "good" and "dead"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        |
      | main   | local, origin | conflicting commit | conflicting_file |
      | dead   | local, origin | dead-end commit    | file             |
      | good   | local, origin | good commit        | file             |
    And the current branch is "good"
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill dead"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | good   | git fetch --prune --tags |
      |        | git push origin :dead    |
      |        | git branch -D dead       |
    And the current branch is still "good"
    And my workspace still contains my uncommitted file
    And the branches are now
      | REPOSITORY    | BRANCHES   |
      | local, origin | main, good |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | conflicting commit |
      | good   | local, origin | good commit        |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | good   | git branch dead {{ sha 'dead-end commit' }} |
      |        | git push -u origin dead                     |
    And the current branch is still "good"
    And my workspace still contains my uncommitted file
    And now the initial commits exist
    And the initial branches and hierarchy exist
