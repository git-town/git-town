Feature: display debug statistics

  Scenario: Git Town command ran successfully
    Given I ran "git-town sync"
    When I run "git-town status --debug"
    Then it prints:
      """
      Ran 7 shell commands.
      """
