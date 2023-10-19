Feature: display all executed Git commands

  Scenario: debug adding aliases
    When I run "git-town aliases add --verbose"
    Then it prints:
      """
      Ran 14 shell commands.
      """

  Scenario: debug removing aliases
    Given I ran "git-town aliases add"
    When I run "git-town aliases remove --verbose"
    Then it prints:
      """
      Ran 14 shell commands.
      """
