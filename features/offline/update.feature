Feature: change offline mode

  Scenario: enable
    When I run "git-town offline true"
    Then offline mode is enabled

  Scenario: disable
    Given Git Town is in offline mode
    When I run "git-town offline false"
    Then offline mode is disabled

  Scenario: invalid value
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """
