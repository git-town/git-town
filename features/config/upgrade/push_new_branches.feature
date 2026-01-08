Feature: automatically upgrade outdated configuration

  Scenario Outline:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And <LOCATION> Git setting "git-town.push-new-branches" is "true"
    And the current branch is "feature"
    When I run "git-town <COMMAND>"
    Then Git Town prints:
      """
      Upgrading deprecated <LOCATION> setting git-town.push-new-branches to git-town.share-new-branches.
      """
    And <LOCATION> Git setting "git-town.new-branch-push-flag" now doesn't exist
    And <LOCATION> Git setting "git-town.share-new-branches" is now "push"

    Examples:
      | COMMAND  | LOCATION |
      | hack foo | local    |
      | hack foo | global   |
