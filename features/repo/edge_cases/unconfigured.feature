@skipWindows
Feature: ask for missing configuration

  Scenario: unconfigured
    Given a Git repo clone
    And Git Town is not configured
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
