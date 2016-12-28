Feature: git town-prepend: errors when trying to prepend something in front of the main branch

  As a developer accidentally trying to prepend someting in front of the main branch
  I should see an error that the main branch has no parents
  So that I know about my mistake and run "git hack" instead.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local and remote | good commit |
    And I am on the "main" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git town-prepend`
    Then it runs no commands
    And I get the error "You can not add a parent for the main branch"
    And I am still on the "main" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I am left with my original commits
