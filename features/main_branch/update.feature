Feature: configure the main branch

  Scenario: not configured
    Given the main branch is not configured
    When I run "git-town main-branch main"
    Then it prints no output
    And the main branch is now "main"

  Scenario: configured
    Given my repo has the branches "old" and "new"
    And the main branch is "old"
    When I run "git-town main-branch new"
    Then it prints no output
    And the main branch is now "new"

  Scenario: non-existing branch
    When I run "git-town main-branch non-existing"
    Then it prints the error:
      """
      no branch named "non-existing"
      """
