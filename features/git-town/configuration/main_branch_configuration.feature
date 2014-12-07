Feature: main-branch configuration

  Scenario: printing the main branch when it's not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch`
    Then I see "Main branch: ''"


  Scenario: printing the main branch when it's configured
    Given I have configured the main branch name as "main"
    When I run `git town main-branch`
    Then I see "Main branch: 'main'"


  Scenario: setting the main branch
    Given I have configured the main branch name as "main-old"
    When I run `git town main-branch main-new`
    Then I see "main branch stored as 'main-new'"
