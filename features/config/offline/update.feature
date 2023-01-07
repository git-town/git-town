Feature: change offline mode

  Scenario: enable
    Given setting "offline" is "false"
    When I run "git-town config offline yes"
    Then setting "offline" is now "true"

  Scenario: disable
    Given setting "offline" is "true"
    When I run "git-town config offline no"
    Then setting "offline" is now "false"

  Scenario: invalid value
    Given setting "offline" is "false"
    When I run "git-town config offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "yes" or "no"
      """
    And setting "offline" is still "false"
