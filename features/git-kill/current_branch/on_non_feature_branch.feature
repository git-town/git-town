Feature: git kill: errors when trying to kill a non-feature branch

  As a developer accidentally trying to kill a non-feature branch
  I should see an error that I cannot delete non-feature branches
  So that my release infrastructure remains intact and my project stays shippable.


  Background:
    Given I have branches named "qa"
    And my non-feature branches are configured as "qa"
    And the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE   | FILE NAME |
      | qa     | local and remote | qa commit | qa_file   |
    And I am on the "qa" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git kill`
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "qa" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main, qa |
      | remote     | main, qa |
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE   | FILE NAME |
      | qa     | local and remote | qa commit | qa_file   |

