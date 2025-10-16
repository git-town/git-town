Feature: enforce stashing via CLI flag

  Background:
    Given a Git repo with origin
    And Git setting "git-town.stash" is "false"
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town hack new --stash"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout -b new         |
      | new    | git stash pop               |
      |        | git restore --staged .      |
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
    And the initial branches and lineage exist now
    And the initial commits exist now
