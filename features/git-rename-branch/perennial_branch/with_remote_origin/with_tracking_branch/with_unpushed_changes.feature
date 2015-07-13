Feature: git rename-branch: errors if renaming a perennial branch that has unpushed changes

  As a developer renaming a perennial branch that has unpushed changes
  I should get an error that the given branch is not in sync with its tracking branch
  So that I know branches must be in sync in order to be renamed.


  Background:
    Given I have a branch named "production"
    And my perennial branches are configured as "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE                  |
      | main       | local and remote | main commit              |
      | production | local and remote | production commit        |
      |            | local            | remote production commit |
    And I am on the "production" branch
    And I have an uncommitted file
    When I run `git rename-branch production renamed-production -f`


  Scenario: result
    Then I get the error "The branch is not in sync with its tracking branch."
    And I get the error "Run 'git sync production' to sync the branch."
    And I end up on the "production" branch
    And I still have my uncommitted file
    And I am left with my original commits
