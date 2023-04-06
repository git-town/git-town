Feature: display debug statistics

  Scenario: debug adding aliases
    When I run "git-town aliases add --debug"
    Then it prints:
      """
      Ran 14 shell commands.
      """

  Scenario: remove adding aliases
    Given I ran "git-town aliases add"
    When I run "git-town aliases add --debug"
    Then it prints:
      """
      Ran 14 shell commands.
      """
