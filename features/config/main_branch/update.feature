Feature: configure the main branch

  Scenario: not configured
    Given local setting "main-branch-name" is ""
    When I run "git-town config main-branch main"
    Then it prints no output
    And local setting "main-branch-name" is now "main"

  Scenario: previously configured
    Given the branches "old" and "new"
    And local setting "main-branch-name" is "old"
    When I run "git-town config main-branch new"
    Then it prints no output
    And local setting "main-branch-name" is now "new"

  Scenario: non-existing branch
    When I run "git-town config main-branch non-existing"
    Then it prints the error:
      """
      there is no branch "non-existing"
      """
