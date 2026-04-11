@skipWindows
Feature: set an empty browser via Git metadata

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And Git setting "git-town.browser" is ""
    And tool "open" is installed
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
