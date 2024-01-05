Feature: display the sync-upstream setting

  Scenario Outline: default setting
    When I run "git-town config sync-upstream <FLAG>"
    Then it prints:
      """
      yes
      """

    Examples:
      | FLAG     |
      |          |
      | --global |

  Scenario Outline: configured in local Git metadata
    Given local Git Town setting "sync-upstream" is "<VALUE>"
    When I run "git-town config sync-upstream"
    Then it prints:
      """
      <OUTPUT>
      """
    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | on    | yes    |
      | true  | yes    |
      | 1     | yes    |
      | t     | yes    |
      | no    | no     |
      | off   | no     |
      | false | no     |
      | f     | no     |
      | 0     | no     |

  Scenario Outline: configured in global Git metadata
    Given global Git Town setting "sync-upstream" is "<VALUE>"
    When I run "git-town config sync-upstream --global"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | no    | no     |

  Scenario Outline: global and local Git metadata set to different values
    Given global Git Town setting "sync-upstream" is "true"
    And local Git Town setting "sync-upstream" is "false"
    When I run "git-town config sync-upstream <FLAG>"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | FLAG     | OUTPUT |
      | --global | yes    |
      |          | no     |

  Scenario: empty config file
    Given the configuration file:
      """
      """
    When I run "git-town config sync-upstream"
    Then it prints:
      """
      yes
      """

  Scenario: set in config file
    Given the configuration file:
      """
      sync-upstream = false
      """
    When I run "git-town config sync-upstream"
    Then it prints:
      """
      no
      """

  Scenario: invalid value in Git config
    Given Git Town setting "sync-upstream" is "zonk"
    When I run "git-town config sync-upstream"
    Then it prints the error:
      """
      invalid value for git-town.sync-upstream: "zonk". Please provide either "yes" or "no"
      """

  Scenario: invalid value in config file
    Given the configuration file:
      """
      sync-upstream = "zonk"
      """
    When I run "git-town config sync-upstream"
    Then it prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content: toml: line 1 (last key "sync-upstream"): incompatible types: TOML value has type string; destination has type boolean
      """
