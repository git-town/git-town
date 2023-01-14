Feature: display the push-hook setting

  Scenario: default local setting
    When I run "git-town config push-hook"
    Then it prints:
      """
      yes
      """

  Scenario: default global setting
    When I run "git-town config push-hook --global"
    Then it prints:
      """
      yes
      """

  Scenario Outline: local setting
    Given setting "push-hook" is "<VALUE>"
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

  @this
  Scenario Outline: global setting
    Given setting "push-hook" is globally "<VALUE>"
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
    Given setting "push-hook" is globally "true"
    When I run "git-town config push-hook"
    Then it prints:
      """
      yes
      """

  Scenario: global and local set
    Given setting "push-hook" is globally "true"
    And setting "push-hook" is "false"
    When I run "git-town config push-hook"
    Then it prints:
      """
      no
      """

  Scenario: invalid value
    Given setting "push-hook" is "zonk"
    When I run "git-town config push-hook"
    Then it prints the error:
      """
      Error: invalid value for git-town.push-hook: "zonk". Please provide either "true" or "false"
      """
