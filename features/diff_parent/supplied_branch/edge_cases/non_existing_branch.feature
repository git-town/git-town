Feature: does not diff non-existing branch

  Scenario: non-existing branch
    When I run "git-town diff-parent non-existing"
    Then it runs no commands
    And it prints the error:
      """
      there is no local branch named "non-existing"
      """
