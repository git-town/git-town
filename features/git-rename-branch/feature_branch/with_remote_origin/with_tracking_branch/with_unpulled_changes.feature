Feature: git rename-branch: errors if renaming a feature branch that has unpulled changes

  As a developer renaming a feature branch that has unpulled changes
  I should get an error that the given branch is not in sync with its tracking branch
  So that I don't lose work by accidentally overwriting existing branches.


  Background:
    Given I have a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE               |
      | main            | local and remote | main commit           |
      | current-feature | local and remote | feature commit        |
      |                 | remote           | remote feature commit |
    And I am on the "current-feature" branch


  Scenario: with open changes
    And I have an uncommitted file
    When I run `git rename-branch current-feature renamed-feature`
    Then I get the error "The branch is not in sync with its tracking branch."
    And I get the error "Run 'git sync current-feature' to sync the branch."
    And I end up on the "current-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits


  Scenario: without open changes
    When I run `git rename-branch current-feature renamed-feature`
    Then I get the error "The branch is not in sync with its tracking branch."
    And I get the error "Run 'git sync current-feature' to sync the branch."
    And I end up on the "current-feature" branch
    And I am left with my original commits
