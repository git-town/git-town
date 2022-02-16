Feature: local repository

  Background:
    Given my repo does not have an origin
    And the local feature branches "good" and "other"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME        |
      | main   | local    | main commit  | conflicting_file |
      | good   | local    | good commit  | file             |
      | other  | local    | other commit | file             |
    And the current branch is "good"
    And an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town kill other"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND             |
      | good   | git branch -D other |
    And the current branch is still "good"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES   |
      | local      | main, good |
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | good   | local    | good commit |
    And this branch hierarchy exists now
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | good   | git branch other {{ sha 'other commit' }} |
    And the current branch is still "good"
    And the uncommitted file still exists
    And now the initial commits exist
    And the initial branches and hierarchy exist
