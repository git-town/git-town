Feature: displaying the new branch push flag configuration


  Scenario: default setting
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """


  Scenario: set to "true"
    Given the new-branch-push-flag configuration is true
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      true
      """


  Scenario: set to "false"
    Given the new-branch-push-flag configuration is false
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      false
      """


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

  Scenario: invalid value
    Given the new-branch-push-flag configuration is "zonk"
    When I run "git-town new-branch-push-flag"
    Then it prints:
      """
      Invalid value for git-town.new-branch-push-flag: "zonk". Please provide either true or false. Considering false for now.
      """
