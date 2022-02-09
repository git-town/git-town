Feature: display the current offline status

  Scenario: default value
    When I run "git-town offline"
    Then it prints:
      """
      false
      """

  Scenario: enabled
    Given the "offline" setting is "true"
    When I run "git-town offline"
    Then it prints:
      """
      true
      """

  Scenario: disabled
    Given the "offline" setting is "false"
    When I run "git-town offline"
    Then it prints:
      """
      false
      """

  Scenario: invalid value
    Given the "offline" setting is "zonk"
    When I run "git-town offline"
    Then it prints:
      """
      Invalid value for git-town.offline: "zonk". Please provide either true or false. Considering false for now.
      """
