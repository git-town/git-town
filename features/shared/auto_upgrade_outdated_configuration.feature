Feature: automatic upgrade outdated configuration

  @this
  Scenario Outline:
    Given local setting "<OLD>" is "<VALUE>"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated local setting "git-town.<OLD>".
      I am upgrading this setting to the new format "git-town.<NEW>".
      """
    And <LOCALITY> setting "<NEW>" is now "<VALUE>"
    And <LOCALITY> setting "<OLD>" no longer exists

    Examples:
      | COMMAND | LOCALITY | OLD                  | NEW               | VALUE |
      | config  | local    | new-branch-push-flag | push-new-branches | true  |
      | config  | global   | new-branch-push-flag | push-new-branches | true  |
