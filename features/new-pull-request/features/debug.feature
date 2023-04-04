@skipWindows
Feature: display debug statistics

  Scenario: debug mode enabled
    Given tool "open" is installed
    And the current branch is a feature branch "feature"
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town new-pull-request --debug"
    Then it prints:
      """
      Ran 31 shell commands.
      """
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
