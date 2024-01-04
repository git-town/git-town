Feature: display the push-new-branches setting

  Scenario Outline: default setting
    When I run "git-town config push-new-branches <FLAG>"
    Then it prints:
      """
      no
      """

    Examples:
      | FLAG     |
      |          |
      | --global |

  Scenario Outline: configured in local Git metadata
    Given local Git Town setting "push-new-branches" is "<VALUE>"
    When I run "git-town config push-new-branches"
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
    Given global Git Town setting "push-new-branches" is "<VALUE>"
    When I run "git-town config push-new-branches --global"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | no    | no     |

  Scenario Outline: global and local set to different values
    Given global Git Town setting "push-new-branches" is "true"
    And local Git Town setting "push-new-branches" is "false"
    When I run "git-town config push-new-branches <FLAG>"
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
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      no
      """

  Scenario: set in config file
    Given the configuration file:
      """
      push-new-branches = true
      """
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      yes
      """

  Scenario: invalid value in Git config
    Given Git Town setting "push-new-branches" is "zonk"
    When I run "git-town config push-new-branches"
    Then it prints the error:
      """
      invalid value for git-town.push-new-branches: "zonk". Please provide either "yes" or "no"
      """

  Scenario: invalid value in config file
    Given the configuration file:
      """
      push-new-branches = zonk
      """
    When I run "git-town config push-new-branches"
    Then it prints the error:
      """
      the configuration file ".git-branches.yml" does not contain TOML-formatted content: toml: line 1 (last key "push-new-branches"): expected value but found "zonk" instead
      """
