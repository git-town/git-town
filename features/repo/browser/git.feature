@skipWindows
Feature: set a custom browser via Git metadata

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And Git setting "git-town.browser" is "firefox"
    And tool "firefox" is installed
    And tool "open" is installed
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                      |
      | main   | frontend | firefox https://github.com/git-town/git-town |
