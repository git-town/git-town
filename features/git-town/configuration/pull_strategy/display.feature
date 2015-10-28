Feature: passing an invalid option to the pull strategy configuration

  As a user or tool configuring Git Town's pull branch strategy
  I want to know what the existing value for the pull-strategy is.


  Scenario: using invalid option
    Given my repository has the "merge" pull strategy configured
    When I run `git town config --pull-strategy`
    Then I see
      """
	  merge
      """
