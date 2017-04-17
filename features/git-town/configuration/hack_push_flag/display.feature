Feature: passing an invalid option to the hack push flag configuration

  As a user or tool configuring Git Town's git-hack push flag
  I want to know what the existing value for the git-hack push flag is
  So I can decide whether to I want to adjust it.


  Scenario: default setting
    When I run `gt hack-push-flag`
    Then I see
      """
      true
      """


  Scenario: set to "true"
    Given my repository has the "hack-push-flag" configuration set to "true"
    When I run `gt hack-push-flag`
    Then I see
      """
      true
      """


  Scenario: set to "false"
    Given my repository has the "hack-push-flag" configuration set to "false"
    When I run `gt hack-push-flag`
    Then I see
      """
      false
      """
