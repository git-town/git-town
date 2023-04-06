Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCATION> setting "new-branch-push-flag" is "true"
    And the current branch is a feature branch "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.new-branch-push-flag".
      I am upgrading this setting to the new format "git-town.push-new-branches".
      """
    And <LOCATION> setting "push-new-branches" is now "true"
    And <LOCATION> setting "new-branch-push-flag" no longer exists

    Examples:
      | COMMAND                  | LOCATION |
      | config                   | local    |
      | config                   | global   |
      | config push-new-branches | local    |
      | config push-new-branches | global   |
      | append foo               | local    |
      | append foo               | global   |
      | hack foo                 | local    |
      | hack foo                 | global   |
      | prepend foo              | local    |
      | prepend foo              | global   |
