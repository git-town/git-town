Feature: prototype another remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE   | PARENT | LOCATIONS |
      | remote-feature | (none) | main   | origin    |
    And I run "git fetch"
    When I run "git-town prototype remote-feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | main   | git checkout remote-feature |
    And Git Town prints:
      """
      branch "remote-feature" is now a prototype branch
      """
    And branch "remote-feature" now has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH         | COMMAND                      |
      | remote-feature | git checkout main            |
      | main           | git branch -D remote-feature |
    And branch "remote-feature" now has type "feature"
    And the initial branches exist now
