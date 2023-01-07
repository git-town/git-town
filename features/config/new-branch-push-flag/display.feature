Feature: display the new-branch-push-flag setting

  Scenario: default local setting
    When I run "git-town config new-branch-push-flag"
    Then it prints:
      """
      no
      """

  Scenario: default global setting
    When I run "git-town config new-branch-push-flag --global"
    Then it prints:
      """
      no
      """

  Scenario Outline: local setting
    Given setting "new-branch-push-flag" is "<VALUE>"
    When I run "git-town config new-branch-push-flag"
    Then it prints:
      """
      <VALUE>
      """
    Examples:
      | VALUE |
      | yes   |
      | no    |

  @this
  Scenario Outline: global setting
    Given setting "new-branch-push-flag" is globally "<VALUE>"
    When I run "git-town config new-branch-push-flag --global"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE |
      | yes   |
      | no    |

  Scenario: global set, local not set
    Given setting "new-branch-push-flag" is globally "true"
    When I run "git-town config new-branch-push-flag"
    Then it prints:
      """
      yes
      """

  Scenario: global and local set
    Given setting "new-branch-push-flag" is globally "true"
    And setting "new-branch-push-flag" is "false"
    When I run "git-town config new-branch-push-flag"
    Then it prints:
      """
      no
      """

  Scenario: invalid value
    Given setting "new-branch-push-flag" is "zonk"
    When I run "git-town config new-branch-push-flag"
    Then it prints:
      """
      Invalid value for git-town.new-branch-push-flag: "zonk". Please provide either "yes" or "no". Considering "no" for now.
      """
