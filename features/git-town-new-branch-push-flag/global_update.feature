Feature: set the global new-branch-push-flag

  Scenario: globally update to "true"
    When I run `git-town new-branch-push-flag --global true`
    Then git is now configured with "new-branch-push-flag" set to "true"


  Scenario: globally update to "false"
    When I run `git-town new-branch-push-flag --global false`
    Then git is now configured with "new-branch-push-flag" set to "false"
