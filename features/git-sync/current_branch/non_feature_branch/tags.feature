Feature: git sync: syncing the current non-feature branch syncs the tags

  As a developer syncing a non-feature branch
  I want my tags to be published
  So that I can do tagging work effectively on my local machine and have more time for other work.


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "production" branch
    And I have the following tags
      | NAME   | LOCATION |
      | local  | local    |
      | remote | remote   |
    When I run `git sync`


  Scenario: result
    Then I now have the following tags
      | NAME   | LOCATION         |
      | local  | local and remote |
      | remote | local and remote |
