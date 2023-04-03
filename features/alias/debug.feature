Feature: display debug statistics

  @this
  Scenario: debug adding aliases
    When I run "git-town aliases add"
    Then it prints:
      """
      Ran 30 shell commands.
      """
