Feature: Git Sync: syncing the main branch pushes tags to the remote



  Background:
    Given I am on the "main" branch
    And I add a local tag "v1.0"
    When I run `git sync`


  Scenario: result
    Then tag "v1.0" has been pushed to the remote
