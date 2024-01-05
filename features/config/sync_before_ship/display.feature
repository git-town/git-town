Feature: display the sync-before-ship setting

  Scenario Outline: default setting
    When I run "git-town config sync-before-ship <FLAG>"
    Then it prints:
      """
      no
      """

    Examples:
      | FLAG     |
      |          |
      | --global |

  Scenario Outline: configured in local Git metadata
    Given local Git Town setting "sync-before-ship" is "<VALUE>"
    When I run "git-town config sync-before-ship"
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
    Given global Git Town setting "sync-before-ship" is "<VALUE>"
    When I run "git-town config sync-before-ship --global"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | no    | no     |

  Scenario Outline: global and local set to different values
    Given global Git Town setting "sync-before-ship" is "true"
    And local Git Town setting "sync-before-ship" is "false"
    When I run "git-town config sync-before-ship <FLAG>"
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
    When I run "git-town config sync-before-ship"
    Then it prints:
      """
      no
      """

  Scenario: set in config file
    Given the configuration file:
      """
      sync-before-ship = true
      """
    When I run "git-town config sync-before-ship"
    Then it prints:
      """
      yes
      """

  Scenario: invalid value in Git config
    Given Git Town setting "sync-before-ship" is "zonk"
    When I run "git-town config sync-before-ship"
    Then it prints the error:
      """
      invalid value for git-town.sync-before-ship: "zonk". Please provide either "yes" or "no"
      """

  Scenario: invalid value in config file
    Given the configuration file:
      """
      sync-before-ship = "zonk"
      """
    When I run "git-town config sync-before-ship"
    Then it prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content: toml: line 1 (last key "sync-before-ship"): incompatible types: TOML value has type string; destination has type boolean
      """
