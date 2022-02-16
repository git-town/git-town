Feature: display the current offline status

  Scenario: default value
    When I run "git-town offline"
    Then it prints:
      """
      false
      """

  Scenario: enabled
    Given setting "offline" is "true"
    When I run "git-town offline"
    Then it prints:
      """
      true
      """

  Scenario: disabled
    Given setting "offline" is "false"
    When I run "git-town offline"
    Then it prints:
      """
      false
      """

  Scenario: invalid value
    Given setting "offline" is "zonk"
    When I run "git-town offline"
    Then it prints:
      """
      Invalid value for git-town.offline: "zonk". Please provide either true or false. Considering false for now.
      """
