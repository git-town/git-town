Feature: enforce stashing via CLI flag

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And Git setting "git-town.stash" is "false"
    And the current branch is "old"
    And an uncommitted file
    When I run "git-town prepend new --stash"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | old    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout -b new main    |
      | new    | git stash pop               |
      |        | git restore --staged .      |
    And this lineage exists now
      """
      main
        new
          old
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | new    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout old            |
      | old    | git branch -D new           |
      |        | git stash pop               |
      |        | git restore --staged .      |
    And the initial branches and lineage exist now
    And the initial commits exist now
