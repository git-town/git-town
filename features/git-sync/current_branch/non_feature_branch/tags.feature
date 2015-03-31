Feature: git sync: syncing the current non-feature branch syncs the tags

  As a developer using Git tags for release management
  I want my tags to be published whenever I sync a non-feature branch
  So that I can do tagging work effectively on my local machine.


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "production" branch
    And I have the following tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | remote-tag | remote   |
    When I run `git sync`


  Scenario: result
    Then I now have the following tags
      | NAME       | LOCATION         |
      | local-tag  | local and remote |
      | remote-tag | local and remote |
