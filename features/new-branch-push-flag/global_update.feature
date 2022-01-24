Feature: set the global new-branch-push-flag

  Scenario: globally update to "true"
    When I run "git-town new-branch-push-flag --global true"
    Then the new-branch-push-flag configuration is now true

  Scenario: globally update to "false"
    When I run "git-town new-branch-push-flag --global false"
    Then the new-branch-push-flag configuration is now false
