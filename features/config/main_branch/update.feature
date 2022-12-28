Feature: configure the main branch

  Scenario: not configured
    Given the main branch is not set
    When I run "git-town config main-branch main"
    Then it prints no output
    And the main branch is now "main"

  Scenario: configured
    Given the branches "old" and "new"
    And the main branch is "old"
    When I run "git-town config main-branch new"
    Then it prints no output
    And the main branch is now "new"

  Scenario: non-existing branch
    When I run "git-town config main-branch non-existing"
    Then it prints the error:
      """
      no branch named "non-existing"
      """
