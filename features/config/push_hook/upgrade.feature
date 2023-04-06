Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCALITY> setting "push-verify" is "true"
    And the current branch is a feature branch "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCALITY> setting "git-town.push-verify".
      I am upgrading this setting to the new format "git-town.push-hook".
      """
    And <LOCALITY> setting "push-hook" is now "true"
    And <LOCALITY> setting "push-verify" no longer exists

    Examples:
      | COMMAND          | LOCALITY |
      | config push-hook | local    |
      | config push-hook | global   |
