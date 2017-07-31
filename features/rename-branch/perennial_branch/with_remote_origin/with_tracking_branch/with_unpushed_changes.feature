Feature: git town-rename-branch: errors if renaming a perennial branch that has unpushed changes

  As a developer renaming a perennial branch that has unpushed changes
  I should get an error that the given branch is not in sync with its tracking branch
  So that I know branches must be in sync in order to be renamed.


  Background:
    Given I have a perennial branch named "production"
    And the following commits exist in my repository
      | BRANCH     | LOCATION         | MESSAGE                  |
      | main       | local and remote | main commit              |
      | production | local and remote | production commit        |
      |            | local            | remote production commit |
    And I am on the "production" branch
    And I have an uncommitted file
    When I run `git-town rename-branch --force production renamed-production`


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND           |
      | production | git fetch --prune |
    And I get the error "'production' is not in sync with its tracking branch. Please sync the branches before renaming."
    And I end up on the "production" branch
    And I still have my uncommitted file
    And I am left with my original commits
