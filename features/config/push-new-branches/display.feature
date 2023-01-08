Feature: display the push-new-branches setting

  Scenario: default local setting
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      no
      """

  Scenario: default global setting
    When I run "git-town config push-new-branches --global"
    Then it prints:
      """
      no
      """

  Scenario Outline: local setting
    Given setting "push-new-branches" is "<VALUE>"
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

  Scenario Outline: global setting
    Given setting "push-new-branches" is globally "<VALUE>"
    When I run "git-town config push-new-branches --global"
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

  Scenario: global set, local not set
    Given setting "push-new-branches" is globally "yes"
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      yes
      """

  Scenario: global and local set
    Given setting "push-new-branches" is globally "yes"
    And setting "push-new-branches" is "false"
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      no
      """

  Scenario: invalid value
    Given setting "push-new-branches" is "zonk"
    When I run "git-town config push-new-branches"
    Then it prints the error:
      """
      Error: invalid value for git-town.push-new-branches: "zonk". Please provide either "yes" or "no"
      """
