Feature: no browser is installed

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And no tool to open browsers is installed
    When I run "git-town repo"

  Scenario: result
    Then Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
