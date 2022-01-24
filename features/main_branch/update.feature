Feature: set the main branch configuration


  Scenario: main branch not yet configured
    Given my repo doesn't have a main branch configured
    When I run "git-town main-branch main"
    Then it prints no output
    And the main branch is now configured as "main"


  Scenario: main branch is configured
    Given my repo has the branches "main-old" and "main-new"
    And the main branch is configured as "main-old"
    When I run "git-town main-branch main-new"
    Then it prints no output
    And the main branch is now configured as "main-new"


  Scenario: invalid branch name
    When I run "git-town main-branch non-existing"
    Then it prints the error:
      """
      no branch named "non-existing"
      """
