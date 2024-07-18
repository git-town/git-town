Feature: dry run appending a new feature branch to an existing feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    And an uncommitted file
    When I run "git-town append new --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND             |
      | existing | git add -A          |
      |          | git stash           |
      |          | git checkout -b new |
      | new      | git stash pop       |
    And the current branch is still "existing"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "existing"
    And the initial commits exist
    And the initial lineage exists
    And the uncommitted file still exists
