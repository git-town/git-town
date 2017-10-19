Feature: passing an invalid option to the hack push flag configuration

  As a user or tool configuring Git Town's git-hack push flag
  I want to know what the existing value for the git-hack push flag is
  So I can decide whether to I want to adjust it.


  Scenario: default setting
    When I run `git-town hack-push-flag`
    Then it prints
      """
      false
      """


  Scenario: set to "true"
    Given the "hack-push-flag" configuration is set to "true"
    When I run `git-town hack-push-flag`
    Then it prints
      """
      true
      """


  Scenario: set to "false"
    Given the "hack-push-flag" configuration is set to "false"
    When I run `git-town hack-push-flag`
    Then it prints
      """
      false
      """
