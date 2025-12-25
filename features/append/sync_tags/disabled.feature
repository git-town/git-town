Feature: don't sync tags while appending

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-tags" is "false"
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    When I run "git-town append new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git fetch --prune --no-tags |
      |        | git checkout -b new         |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout main |
      | main   | git branch -D new |
    And the initial tags exist now
