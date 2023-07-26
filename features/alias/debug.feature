Feature: display debug statistics

  Scenario: debug adding aliases
    When I run "git-town aliases add --debug"
    Then it prints:
      """
      Ran 15 shell commands.
      """

  Scenario: debug removing aliases
    Given I ran "git-town aliases add"
    When I run "git-town aliases remove --debug"
    Then it prints:
      """
      Ran 15 shell commands.
      """
