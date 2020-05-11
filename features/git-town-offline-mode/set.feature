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
