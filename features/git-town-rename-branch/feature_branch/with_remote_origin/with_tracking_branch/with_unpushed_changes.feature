Feature: git town-rename-branch: errors if renaming a feature branch that has unpushed changes

  As a developer renaming a feature branch that has unpushed changes
  I should get an error that the given branch is not in sync with its tracking branch
  So that I know branches must be in sync in order to be renamed.


  Background:
    Given I have a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE              |
      | main            | local and remote | main commit          |
      | current-feature | local and remote | feature commit       |
      |                 | local            | local feature commit |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git-town rename-branch current-feature renamed-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND           |
      | current-feature | git fetch --prune |
    And I get the error "'current-feature' is not in sync with its tracking branch. Please sync the branches before renaming."
    And I end up on the "current-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
