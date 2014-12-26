Feature: display the main branch configuration

  As a user or tool unsure about which branch is currently configured as the main development branch
  I want to be able to see this information simply and directly
  So that I can use it without furter thinking or processing, and my Git Town workflows are effective.


  Scenario: main branch not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch`
    Then I see "Main branch: [none]"


  Scenario: main branch is configured
    Given I have configured the main branch name as "main"
    When I run `git town main-branch`
    Then I see "Main branch: main"
