Feature: git sync: syncing a feature branch pulls tags

  As a developer using tags for release management

  I want that tags are pulled automatically for me whenever I sync
  So that my local workspace has the same tags that exist on the remote


  Background:
    Given I have a feature branch named "feature"
    And I have the following tags
      | NAME   | LOCATION |
      | local  | local    |
      | remote | remote   |
    And I am on the "feature" branch
    And I run `git sync`


  Scenario: result
    Then I now have the following tags
      | NAME   | LOCATION         |
      | local  | local            |
      | remote | local and remote |
