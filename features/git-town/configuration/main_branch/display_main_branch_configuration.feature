Feature: display the main branch configuration

  Scenario: main branch not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch` while allowing errors
    Then I see "[none]"


  Scenario: main branch is configured
    Given I have configured the main branch name as "main"
    When I run `git town main-branch`
    Then I see "main"
