Feature: make the current local feature branch an observed branch

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town observe feature"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      branch "feature" is local only - branches you want to observe must have a remote branch because they are per definition other people's branches
      """
    And branch "feature" is still a feature branch
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "feature" is still a feature branch
    And there are still no observed branches
