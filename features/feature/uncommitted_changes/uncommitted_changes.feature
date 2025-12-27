Feature: convert a branch to feature branch in the presence of uncommitted changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | PARENT | LOCATIONS |
      | existing | parked | main   | local     |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town feature existing"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch existing is now a feature branch
      """
    And the initial branches and lineage exist now
    And branch "existing" now has type "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git stash pop               |
      |        | git restore --staged .      |
    And the initial branches and lineage exist now
    And branch "existing" now has type "parked"
    And the initial commits exist now
