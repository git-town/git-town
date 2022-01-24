Feature: git new-pull-request: offline mode

  When in offline mode
  I want to be told that Git Town is in offline mode
  So that I know I cannot create new pull requests.

  Scenario: trying to create a new pull request in offline mode
    Given Git Town is in offline mode
    When I run "git-town new-pull-request"
    Then it prints the error:
      """
      this command requires an active internet connection
      """
