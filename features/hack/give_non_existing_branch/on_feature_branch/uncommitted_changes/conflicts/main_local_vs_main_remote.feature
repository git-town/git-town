Feature: conflicts between the main branch and its tracking branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git add -A               |
      |          | git stash                |
      |          | git checkout -b new main |
      | new      | git stash pop            |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE                   |
      | main   | local    | conflicting local commit  |
      |        | origin   | conflicting origin commit |
      | new    | local    | conflicting local commit  |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial branches and lineage exist
