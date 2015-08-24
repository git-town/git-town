Feature: git kill: errors when trying to kill a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given I have branches named "feature" and "qa"
    And my perennial branches are configured as "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git kill qa`
    Then it runs no commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
