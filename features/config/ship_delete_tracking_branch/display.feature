Feature: display the ship-delete-tracking-branch setting

  Scenario Outline: default setting in Git metadata
    When I run "git-town config ship-delete-tracking-branch <FLAG>"
    Then it prints:
      """
      yes
      """

    Examples:
      | FLAG     |
      |          |
      | --global |

  Scenario Outline: configured in local Git metadata
    Given local Git Town setting "ship-delete-tracking-branch" is "<VALUE>"
    When I run "git-town config ship-delete-tracking-branch"
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
    Given global Git Town setting "ship-delete-tracking-branch" is "<VALUE>"
    When I run "git-town config ship-delete-tracking-branch --global"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | no    | no     |

  Scenario Outline: global and local set to different values
    Given global Git Town setting "ship-delete-tracking-branch" is "true"
    And local Git Town setting "ship-delete-tracking-branch" is "false"
    When I run "git-town config ship-delete-tracking-branch <FLAG>"
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
    When I run "git-town config ship-delete-tracking-branch"
    Then it prints:
      """
      yes
      """

  Scenario: set in config file
    Given the configuration file:
      """
      ship-delete-tracking-branch = false
      """
    When I run "git-town config ship-delete-tracking-branch"
    Then it prints:
      """
      no
      """

  Scenario: invalid value in Git config
    Given Git Town setting "ship-delete-tracking-branch" is "zonk"
    When I run "git-town config ship-delete-tracking-branch"
    Then it prints the error:
      """
      invalid value for git-town.ship-delete-tracking-branch: "zonk". Please provide either "yes" or "no"
      """

  Scenario: invalid value in config file
    Given the configuration file:
      """
      ship-delete-tracking-branch = "zonk"
      """
    When I run "git-town config ship-delete-tracking-branch"
    Then it prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content: toml: line 1 (last key "ship-delete-tracking-branch"): incompatible types: TOML value has type string; destination has type boolean
      """
