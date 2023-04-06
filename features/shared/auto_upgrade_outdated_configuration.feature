Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given <LOCALITY> setting "<OLD>" is "<VALUE>"
    Given the current branch is a feature branch "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCALITY> setting "git-town.<OLD>".
      I am upgrading this setting to the new format "git-town.<NEW>".
      """
    And <LOCALITY> setting "<NEW>" is now "<VALUE>"
    And <LOCALITY> setting "<OLD>" no longer exists

    Examples:
      | COMMAND                  | LOCALITY | OLD                  | NEW               | VALUE |
      | config                   | local    | new-branch-push-flag | push-new-branches | true  |
      | config                   | global   | new-branch-push-flag | push-new-branches | true  |
      | config push-new-branches | local    | new-branch-push-flag | push-new-branches | true  |
      | config push-new-branches | global   | new-branch-push-flag | push-new-branches | true  |
      | append foo               | local    | new-branch-push-flag | push-new-branches | true  |
      | append foo               | global   | new-branch-push-flag | push-new-branches | true  |
      | hack foo                 | local    | new-branch-push-flag | push-new-branches | true  |
      | hack foo                 | global   | new-branch-push-flag | push-new-branches | true  |
      | prepend foo              | local    | new-branch-push-flag | push-new-branches | true  |
      | prepen d foo             | global   | new-branch-push-flag | push-new-branches | true  |
