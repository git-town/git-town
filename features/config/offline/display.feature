Feature: display the current offline status

  Scenario: default value
    When I run "git-town config offline"
    Then it prints:
      """
      no
      """

  Scenario: enabled
    Given setting "offline" is "true"
    When I run "git-town config offline"
    Then it prints:
      """
      yes
      """

  Scenario: disabled
    Given setting "offline" is "false"
    When I run "git-town config offline"
    Then it prints:
      """
      no
      """

  Scenario: invalid value
    Given setting "offline" is "zonk"
    When I run "git-town config offline"
    Then it prints:
      """
      Invalid value for git-town.offline: "zonk". Please provide either "yes" or "no". Considering "no" for now.
      """
