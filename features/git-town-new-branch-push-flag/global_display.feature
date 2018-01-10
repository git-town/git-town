Feature: displaying the global new branch push flag configuration

  Scenario: default set
    When I run `git-town new-branch-push-flag --global`
    Then it prints
      """
      false
      """


  Scenario: set to "true"
    Given the global "new-branch-push-flag" configuration is set to "true"
    When I run `git-town new-branch-push-flag --global`
    Then it prints
      """
      true
      """


  Scenario: set to false
    Given the global "new-branch-push-flag" configuration is set to "false"
    When I run `git-town new-branch-push-flag --global`
    Then it prints
      """
      false
      """
