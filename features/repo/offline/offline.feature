Feature: offline mode

  Scenario: try to prune branches in offline mode
    Given a Git repo with origin
    And offline mode is enabled
    When I run "git-town repo"
    Then Git Town prints the error:
      """
      this command requires an active internet connection
      """
