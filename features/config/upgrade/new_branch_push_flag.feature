Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And <LOCATION> Git Town setting "new-branch-push-flag" is "true"
    And the current branch is "feature"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      Upgrading deprecated <LOCATION> setting "git-town.new-branch-push-flag" to "git-town.push-new-branches".
      """
    And <LOCATION> Git Town setting "push-new-branches" is now "true"
    And <LOCATION> Git Town setting "new-branch-push-flag" now doesn't exist

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
