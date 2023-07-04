Feature: automatically upgrade outdated configuration

  @this
  Scenario Outline:
    Given <LOCATION> setting "push-verify" is "true"
    And the current branch is a feature branch "feature"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town <COMMAND>"
    Then it prints:
      """
      I found the deprecated <LOCATION> setting "git-town.push-verify".
      I am upgrading this setting to the new format "git-town.push-hook".
      """
    And <LOCATION> setting "push-hook" is now "true"
    And <LOCATION> setting "push-verify" no longer exists

    Examples:
      | COMMAND | LOCATION |
# | config            | local    |
# | config            | global   |
# | config push-hook  | local    |
# | config push-hook  | global   |
# | append foo        | local    |
# | append foo        | global   |
# | hack foo          | local    |
# | hack foo          | global   |
# | prepend foo       | local    |
# | prepend foo       | global   |
# | sync    | local    |
# | sync    | global   |
# | kill              | local    |
# | kill              | global   |
# | new-pull-request  | local    |
# | new-pull-request  | global   |
# | rename-branch bar | local    |
# | rename-branch bar | global   |
