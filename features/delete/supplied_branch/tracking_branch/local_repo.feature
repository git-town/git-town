Feature: local repository

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | good  | feature | main   | local     |
      | other | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME |
      | good   | local    | good commit  | file      |
      | other  | local    | other commit | file      |
    And the current branch is "good"
    When I run "git-town delete other"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | good   | git branch -D other |
    And the current branch is still "good"
    And the branches are now
      | REPOSITORY | BRANCHES   |
      | local      | main, good |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE     |
      | good   | local    | good commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | good   | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | good   | git branch other {{ sha 'other commit' }} |
    And the current branch is still "good"
    And the initial commits exist now
    And the initial branches and lineage exist now
