Feature: change offline mode

  Scenario: enable offline mode
    When I run "git-town offline true"
    Then offline mode is enabled

  Scenario: disable offline mode
    Given Git Town is in offline mode
    When I run "git-town offline false"
    Then offline mode is disabled

  Scenario: invalid configuration setting
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """

  Scenario: multiple values
    When I run "git-town offline true false"
    Then it prints the error:
      """
      accepts at most 1 arg(s), received 2
      """
