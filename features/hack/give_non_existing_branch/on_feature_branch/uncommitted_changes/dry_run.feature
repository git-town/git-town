Feature: dry-run hacking a new feature branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And an uncommitted file
    When I run "git-town hack new --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git add -A               |
      |          | git stash                |
      |          | git checkout -b new main |
      | new      | git stash pop            |
    And the current branch is still "existing"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "existing"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
