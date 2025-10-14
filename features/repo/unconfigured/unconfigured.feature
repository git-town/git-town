@skipWindows
Feature: ask for missing configuration

  Scenario: unconfigured
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And Git Town is not configured
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | main   | open https://github.com/git-town/git-town |
