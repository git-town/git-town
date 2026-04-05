@skipWindows
Feature: set a custom browser via the CLI

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "firefox" is installed
    When I run "git-town repo --browser=firefox"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                      |
      | main   | frontend | firefox https://github.com/git-town/git-town |
