Feature: set the main branch configuration

  As a user or tool configuring Git Town
  I want an easy way to specifically set the main branch
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: main branch not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch main`
    Then I see no output
    And the main branch name is now configured as "main"


  Scenario: main branch is configured
    Given I have branches named "main-old" and "main-new"
    And I have configured the main branch name as "main-old"
    When I run `git town main-branch main-new`
    Then I see no output
    And the main branch name is now configured as "main-new"


  Scenario: invalid branch name
    When I run `git town main-branch non-existing` it errors
    Then I see
      """
      error: no branch named 'non-existing'
      """
