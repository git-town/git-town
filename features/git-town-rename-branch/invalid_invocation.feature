Feature: git town-rename-branch: requires 1 or 2 branch names, and an optional force flag

  As a developer invoking town-rename-branch with incorrect arity
  I should be reminded that I have to provide the branch names to this command
  So that I can use it correctly without having to look that fact up in the readme.


  Background:
    Given I have a feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION | MESSAGE        |
      | current-feature | local    | feature commit |
    And I am on the "current-feature" branch


  Scenario: no branch names given
    When I run `gt rename-branch`
    Then it runs no commands
    And I get the error "Too few arguments"
    And I am still on the "current-feature" branch
    And I am left with my original commits


  Scenario: three branch names given
    When I run `gt rename-branch one two three`
    Then it runs no commands
    And I get the error "Too many arguments"
    And I am still on the "current-feature" branch
    And I am left with my original commits
