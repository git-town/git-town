Feature: display the push-hook setting

  Scenario Outline: default settings
    When I run "git-town config push-hook <SWITCH>"
    Then it prints:
      """
      yes
      """

    Examples:
      | SWITCH   |
      | --global |
      |          |

  Scenario Outline: configured locally
    Given local setting "push-hook" is "<VALUE>"
    When I run "git-town config push-hook"
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

  Scenario Outline: configured globally
    Given global setting "push-hook" is "<VALUE>"
    When I run "git-town config push-hook --global"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | on    | yes    |
      | true  | yes    |
      | t     | yes    |
      | 1     | yes    |
      | no    | no     |
      | off   | no     |
      | false | no     |
      | f     | no     |
      | 0     | no     |

  Scenario: global set, local not set
    Given global setting "push-hook" is "true"
    When I run "git-town config push-hook"
    Then it prints:
      """
      yes
      """

  Scenario Outline: global and local set to different values
    Given global setting "push-hook" is "true"
    And local setting "push-hook" is "false"
    When I run "git-town config push-hook <FLAG>"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | FLAG     | OUTPUT |
      | --global | yes    |
      |          | no     |

  Scenario: invalid value
    Given local setting "push-hook" is "zonk"
    When I run "git-town config push-hook"
    Then it prints the error:
      """
      Error: invalid value for git-town.push-hook: "zonk". Please provide either "yes" or "no"
      """
