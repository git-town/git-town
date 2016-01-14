Feature: passing an invalid option to the pull strategy configuration

  As a user or tool configuring Git Town's pull branch strategy
  I want to know what the existing value for the pull-strategy is.

  Scenario: default
    When I run `git town pull-branch-strategy`
    Then I see
      """
      rebase
      """


  Scenario: explicit rebase
    Given my repository has the "rebase" pull branch strategy configured
    When I run `git town pull-branch-strategy`
    Then I see
      """
      rebase
      """

  Scenario: explicit merge
    Given my repository has the "merge" pull branch strategy configured
    When I run `git town pull-branch-strategy`
    Then I see
      """
      merge
      """
