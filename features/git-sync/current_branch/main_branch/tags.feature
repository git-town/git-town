Feature: git sync: syncing the main branch pushes tags to the remote

  As a developer syncing the main branch
  I want my tags to be published
  So that tags are shared with the team


  Background:
    Given I am on the "main" branch
    And I add a local tag "v1.0"
    When I run `git sync`


  Scenario: result
    Then tag "v1.0" has been pushed to the remote
