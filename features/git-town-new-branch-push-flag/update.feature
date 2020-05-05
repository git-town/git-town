Feature: set the new-branch-push-flag

  As a user or tool configuring Git Town
  I want an easy way to specifically set the new branch push flag
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: update to "true"
    When I run "git-town new-branch-push-flag true"
    Then the new-branch-push-flag configuration is now true


  Scenario: update to "false"
    When I run "git-town new-branch-push-flag false"
    Then the new-branch-push-flag configuration is now false
