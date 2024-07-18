Feature: local repository

  Background:
    Given a local Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | good  | feature | main   | local     |
      | other | feature | main   | local     |
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
      | good   | git add -A          |
      |        | git stash           |
      |        | git branch -D other |
      |        | git stash pop       |
    And the current branch is still "good"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY | BRANCHES   |
      | local      | main, good |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | good   | local    | good commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                   |
      | good   | git add -A                                |
      |        | git stash                                 |
      |        | git branch other {{ sha 'other commit' }} |
      |        | git stash pop                             |
    And the current branch is still "good"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist
