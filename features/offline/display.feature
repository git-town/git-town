Feature: display the current offline status

  Scenario: default value
    When I run "git-town offline"
    Then it prints:
      """
      false
      """

  Scenario: enabled
    Given Git Town is in offline mode
    When I run "git-town offline"
    Then it prints:
      """
      true
      """

  Scenario: invalid value
    Given Git Town's "offline" setting is "zonk"
    When I run "git-town offline"
    Then it prints:
      """
      Invalid value for git-town.offline: "zonk". Please provide either true or false. Considering false for now.
      """
