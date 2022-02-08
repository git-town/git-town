Feature: change offline mode

  Scenario: enable
    When I run "git-town offline true"
    Then Git Town's "offline" setting is now "true"
    And Git Town is now in offline mode

  Scenario: disable
    Given Git Town is in offline mode
    When I run "git-town offline false"
    Then Git Town's "offline" setting is now "false"
    And Git Town is no longer in offline mode

  Scenario: invalid value
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """
    And Git Town is no longer in offline mode
