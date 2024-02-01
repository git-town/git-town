Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> Git Town setting "new-branch-push-flag" is "true"
    And the current branch is a feature branch "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.new-branch-push-flag".
      I am upgrading this setting to the new format "git-town.push-new-branches".
      """
    And <LOCATION> Git Town setting "push-new-branches" is now "true"
    And <LOCATION> Git Town setting "new-branch-push-flag" now doesn't exist

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
