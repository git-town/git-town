Feature: automatically upgrade outdated configuration

  @debug @this
  Scenario Outline:
    Given <LOCALITY> setting "new-branch-push-flag" is "true"
    And the current branch is a feature branch "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCALITY> setting "git-town.new-branch-push-flag".
      I am upgrading this setting to the new format "git-town.push-new-branches".
      """
    And <LOCALITY> setting "push-new-branches" is now "true"
    And <LOCALITY> setting "new-branch-push-flag" no longer exists

    Examples:
      | COMMAND                  | LOCALITY |
      | append foo               | local    |
      | append foo               | global   |
      | config                   | local    |
      | config                   | global   |
      | config push-new-branches | local    |
      | config push-new-branches | global   |
      | hack foo                 | local    |
      | hack foo                 | global   |
      | prepend foo              | local    |
      | prepend foo              | global   |
