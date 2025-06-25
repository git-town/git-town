@skipWindows
Feature: print the URL when the browser crashes

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And tool "open" is broken
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | main   | open https://github.com/git-town/git-town |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
