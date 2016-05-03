Feature: passing an invalid option to the pull strategy configuration

  As a user or tool configuring Git Town's pull branch strategy
  I want to know what the existing value for the pull-strategy is
  So I can decide whether to I want to adjust it.


  Scenario: default setting
    When I run `git town pull-branch-strategy`
    Then I see
      """
      rebase
      """


  Scenario: explicit rebase
    Given my repository has the "pull-branch-strategy" configuration set to "rebase"
    When I run `git town pull-branch-strategy`
    Then I see
      """
      rebase
      """


  Scenario: explicit merge
    Given my repository has the "pull-branch-strategy" configuration set to "merge"
    When I run `git town pull-branch-strategy`
    Then I see
      """
      merge
      """
