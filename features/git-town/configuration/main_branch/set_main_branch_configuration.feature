Feature: set the main branch configuration

  Scenario: main branch not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch main`
    Then I see "main branch stored as 'main'"


  Scenario: main branch is configured
    Given I have configured the main branch name as "main-old"
    When I run `git town main-branch main-new`
    Then I see "main branch stored as 'main-new'"
