Feature: on a forked repo

  Background:
    Given an upstream repo
    And the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND             |
      | main   | git add -A          |
      |        | git stash           |
      |        | git branch new main |
      |        | git checkout new    |
      | new    | git stash pop       |
    And the current branch is now "new"
    And the uncommitted file still exists
    And the initial commits exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git add -A        |
      |        | git stash         |
      |        | git checkout main |
      | main   | git branch -D new |
      |        | git stash pop     |
    And the current branch is now "main"
    And the initial commits exist
    And no lineage exists now
    And the uncommitted file still exists
