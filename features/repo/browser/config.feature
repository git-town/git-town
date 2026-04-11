@skipWindows
Feature: set a custom browser via the config file

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And the committed configuration file:
      """
      [hosting]
      browser = "firefox"
      """
    And tool "firefox" is installed
    And tool "open" is installed
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                      |
      | main   | frontend | firefox https://github.com/git-town/git-town |
