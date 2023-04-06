Feature: display debug statistics

  Scenario:
    Given the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    When I run "git-town repo --debug"
    Then it prints:
      """
      Ran 13 shell commands.
      """
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """
