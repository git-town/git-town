Feature: configure the main branch

  Scenario: not configured
    Given my repo doesn't have a main branch configured
    When I run "git-town main-branch main"
    Then it prints no output
    And the main branch is now "main"

  Scenario: configured
    Given my repo has the branches "main-old" and "main-new"
    And the main branch is "main-old"
    When I run "git-town main-branch main-new"
    Then it prints no output
    And the main branch is now "main-new"

  Scenario: non-existing branch
    When I run "git-town main-branch non-existing"
    Then it prints the error:
      """
      no branch named "non-existing"
      """
