Feature: git town-rename-branch: rename current branch implicitly

  As a developer wishing to rename the current branch
  I should be able reference the current branch implicitly
  So that I can perform my rename quickly.


  Background:
    Given I have a feature branch named "feature"
    And I have a perennial branch named "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE     |
      | main       | local and remote | main commit |
      | production | local and remote | main commit |
      | feature    | local and remote | main commit |


  Scenario: rename feature branch
    When I am on the "feature" branch
    And I run `git town-rename-branch renamed-feature`
    Then I end up on the "renamed-feature" branch
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE     |
      | main            | local and remote | main commit |
      | production      | local and remote | main commit |
      | renamed-feature | local and remote | main commit |


  Scenario: rename perennial branch
    When I am on the "production" branch
    And I run `git town-rename-branch renamed-production -f`
    Then I end up on the "renamed-production" branch
    And I have the following commits
      | BRANCH             | LOCATION         | MESSAGE     |
      | main               | local and remote | main commit |
      | feature            | local and remote | main commit |
      | renamed-production | local and remote | main commit |
