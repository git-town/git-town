Feature: automatically upgrade outdated configuration

  @debug @this
  Scenario Outline:
    Given <LOCALITY> setting "<OLD>" is "<VALUE>"
    And the current branch is a feature branch "feature"
    When I run "git-town <COMMAND> --debug"
    Then it prints:
      """
      I found the deprecated <LOCALITY> setting "git-town.<OLD>".
      I am upgrading this setting to the new format "git-town.<NEW>".
      """
    And <LOCALITY> setting "<NEW>" is now "<VALUE>"
    And <LOCALITY> setting "<OLD>" no longer exists

    Examples:
      | COMMAND | LOCALITY | OLD | NEW | VALUE |
# | config push-hook | local    | push-verify | push-hook | true  |
# | config push-hook | global   | push-verify | push-hook | true  |
