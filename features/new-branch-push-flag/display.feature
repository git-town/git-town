Feature: display the new-branch-push-flag setting

  Scenario: display the default local setting
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """

  Scenario Outline: display the local setting
    Given the new-branch-push-flag configuration is <VALUE>
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      <VALUE>
      """
    Examples:
      | VALUE |
      | true  |
      | false |

  Scenario: globally set to "true", local unset
    Given the global new-branch-push-flag configuration is true
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      true
      """

  Scenario: globally set to "true", local set to "false"
    Given the global new-branch-push-flag configuration is true
    And the new-branch-push-flag configuration is false
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """

  Scenario: invalid configuration setting
    Given the new-branch-push-flag configuration is "zonk"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      Invalid value for git-town.new-branch-push-flag: "zonk". Please provide either true or false. Considering false for now.
      """

  Scenario: display the default global value
    When I run "git-town new-branch-push-flag --global"
    Then it prints:
      """
      false
      """

  Scenario Outline: display global value
    Given the global new-branch-push-flag configuration is <VALUE>
    When I run "git-town new-branch-push-flag --global"
    Then it prints:
      """
      <VALUE>
      """

    Examples:
      | VALUE |
      | true  |
      | false |
