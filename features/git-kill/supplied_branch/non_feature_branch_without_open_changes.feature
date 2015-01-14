Feature: git kill: don't remove a given non-feature branch (without open changes)

  (see ../current_branch/non_feature_branch_with_open_changes.feature)


  Background:
    Given I have branches named "feature" and "qa"
    And my non-feature branches are configured as "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
    And I am on the "feature" branch
    When I run `git kill qa` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | qa      | local and remote | qa commit   | qa_file   |
      | feature | local and remote | good commit | good_file |
