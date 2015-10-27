Feature: passing an invalid option to the pull strategy configuration

  As a user or tool configuring Git Town's pull branch strategy
  I want to know when I use an invalid option
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: using invalid option
    When I run `git town config --pull-strategy woof`
    Then I get
      """
	  usage: git town config --pull-strategy [(merge | rebase)]
      """
