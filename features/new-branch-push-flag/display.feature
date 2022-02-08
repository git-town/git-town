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
    Given the new-branch-push-flag configuration is "<VALUE>"
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
    Given the global new-branch-push-flag configuration is "<VALUE>"
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
    Given the global new-branch-push-flag configuration is "true"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      true
      """

  Scenario: global and local set
    Given the global new-branch-push-flag configuration is "true"
    And the new-branch-push-flag configuration is "false"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """

  Scenario: invalid value
    Given the new-branch-push-flag configuration is "zonk"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      Invalid value for git-town.new-branch-push-flag: "zonk". Please provide either true or false. Considering false for now.
      """
