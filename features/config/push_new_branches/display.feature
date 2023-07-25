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

  Scenario Outline: configured locally
    Given local setting "push-new-branches" is "<VALUE>"
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

  Scenario Outline: configured globally
    Given global setting "push-new-branches" is "<VALUE>"
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
    Given global setting "push-new-branches" is "true"
    When I run "git-town config push-new-branches"
    Then it prints:
      """
      yes
      """

  Scenario Outline: global and local set to different values
    Given global setting "push-new-branches" is "true"
    And local setting "push-new-branches" is "false"
    When I run "git-town config push-new-branches <FLAG>"
    Then it prints:
      """
      <OUTPUT>
      """

    Examples:
      | FLAG     | OUTPUT |
      | --global | yes    |
      |          | no     |

  Scenario: invalid value
    Given setting "push-new-branches" is "zonk"
    When I run "git-town config push-new-branches"
    Then it prints the error:
      """
      Error: invalid value for git-town.push-new-branches: "zonk". Please provide either "true" or "false"
      """
