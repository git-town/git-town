Feature: park another remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE    | PARENT | LOCATIONS |
      | remote-feature | feature | main   | origin    |
    And I run "git fetch"
    When I run "git-town park remote-feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      |        | git checkout remote-feature |
    And Git Town prints:
      """
      branch "remote-feature" is now parked
      """
    And branch "remote-feature" now has type "parked"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH         | COMMAND                      |
      | remote-feature | git checkout main            |
      | main           | git branch -D remote-feature |
    And branch "remote-feature" now has type "feature"
