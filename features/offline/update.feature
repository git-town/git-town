Feature: change offline mode

  Scenario: enable
    Given the "offline" setting is "false"
    When I run "git-town offline true"
    Then Git Town's "offline" setting is now "true"

  Scenario: disable
    Given the "offline" setting is "true"
    When I run "git-town offline false"
    Then Git Town's "offline" setting is now "false"

  Scenario: invalid value
    Given the "offline" setting is "false"
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """
    And Git Town's "offline" setting is still "false"
