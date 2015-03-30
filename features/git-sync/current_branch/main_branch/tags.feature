Feature: git sync: syncing the main branch syncs the tags

  As a developer using Git tags for release management
  I want my tags to be published whenever I sync my main branch
  So that I can do tagging work effectively on my local machine.


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
