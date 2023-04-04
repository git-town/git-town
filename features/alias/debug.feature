Feature: display debug statistics

  @this
  Scenario: debug adding aliases
    When I run "git-town aliases add --debug"
    Then it prints:
      """
      Ran 14 shell commands.
      """
