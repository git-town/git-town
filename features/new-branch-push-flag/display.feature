Feature: display the new-branch-push-flag setting

  Scenario: default local setting
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """

  Scenario: default global setting
    When I run "git-town new-branch-push-flag --global"
    Then it prints:
      """
      false
      """

  Scenario Outline: local setting
    Given setting "new-branch-push-flag" is "<VALUE>"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      <VALUE>
      """
    Examples:
      | VALUE |
      | true  |
      | false |

  Scenario Outline: global setting
    Given setting "new-branch-push-flag" is globally "<VALUE>"
    When I run "git-town new-branch-push-flag --global"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE |
      | true  |
      | false |

  Scenario: global set, local not set
    Given setting "new-branch-push-flag" is globally "true"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      true
      """

  Scenario: global and local set
    Given setting "new-branch-push-flag" is globally "true"
    And setting "new-branch-push-flag" is "false"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """

  Scenario: invalid value
    Given setting "new-branch-push-flag" is "zonk"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      Invalid value for git-town.new-branch-push-flag: "zonk". Please provide either true or false. Considering false for now.
      """
