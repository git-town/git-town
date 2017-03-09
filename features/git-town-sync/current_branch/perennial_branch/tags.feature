Feature: gt sync: syncing the current perennial branch syncs the tags

  As a developer using Git tags for release management
  I want my tags to be published whenever I sync a perennial branch
  So that I can do tagging work effectively on my local machine.


  Background:
    Given I have perennial branches named "production" and "qa"
    And I am on the "production" branch
    And I have the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    When I run `gt sync`


  Scenario: result
    Then I now have the following tags
      | NAME       | LOCATION         |
      | local-tag  | local and remote |
      | remote-tag | local and remote |
