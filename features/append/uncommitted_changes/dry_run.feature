Feature: dry run appending a new feature branch to an existing feature branch

  Background:
    Given a Git repo with origin
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
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout -b new         |
      | new      | git stash pop               |
      |          | git restore --staged .      |
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial lineage exists now
    And the uncommitted file still exists
