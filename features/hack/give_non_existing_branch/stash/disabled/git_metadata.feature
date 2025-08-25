Feature: permanently disable stashing via Git metadata

  Background:
    Given a Git repo with origin
    And the current branch is "main"
    And an uncommitted file
    And Git setting "git-town.stash" is "false"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | main   | git checkout -b new |
    And this lineage exists now
      """
      main
        new
      """
  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | new    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout main           |
      | main   | git branch -D new           |
      |        | git stash pop               |
      |        | git restore --staged .      |
    And the initial commits exist now
    And the initial branches and lineage exist now
