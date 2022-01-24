Feature: enabling offline mode

  When developing on an airplane
  I want to be able to use Git Town without interactions with remote origins
  So that I can work on my code even without internet connection.

  Scenario: enabling offline mode
    When I run "git-town offline true"
    Then offline mode is enabled

  Scenario: disabling offline mode
    Given Git Town is in offline mode
    When I run "git-town offline false"
    Then offline mode is disabled

  Scenario: invalid value
    When I run "git-town offline zonk"
    Then it prints the error:
      """
      invalid argument: "zonk". Please provide either "true" or "false"
      """

  Scenario: multiple values
    When I run "git-town offline true false"
    Then it prints the error:
      """
      accepts at most 1 arg(s), received 2
      """
