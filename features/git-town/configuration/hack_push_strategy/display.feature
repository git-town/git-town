Feature: passing an invalid option to the hack push strategy configuration

  As a user or tool configuring Git Town's git-hack push strategy
  I want to know what the existing value for the git-hack push strategy is
  So I can decide whether to I want to adjust it.


  Scenario: default setting
    When I run `git town hack-push-strategy`
    Then I see
      """
      push
      """


  Scenario: explicit push
    Given my repository has the "hack-push-strategy" configuration set to "push"
    When I run `git town hack-push-strategy`
    Then I see
      """
      push
      """


  Scenario: explicit local
    Given my repository has the "hack-push-strategy" configuration set to "local"
    When I run `git town hack-push-strategy`
    Then I see
      """
      local
      """
