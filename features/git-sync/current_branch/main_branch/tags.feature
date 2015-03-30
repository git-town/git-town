Feature: git sync: syncing the main branch syncs the tags

  As a developer syncing the main branch
  I want my tags to be published
  So that tags are shared with the team


  Background:
    Given I have the following tags
      | NAME   | LOCATION |
      | local  | local    |
      | remote | remote   |
    And I am on the "main" branch
    When I run `git sync`


  Scenario: result
    Then I now have the following tags
      | NAME   | LOCATION         |
      | local  | local and remote |
      | remote | local and remote |
