Feature: syncing the main branch pushes tags to the remote


  Scenario: Tags
    Given I am on the main branch
    And I add a local tag "v1.0"
    When I run `git sync`
    Then tag "v1.0" has been pushed to the remote
