@skipWindows
Feature: disable the browser via the CLI

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --no-browser"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
