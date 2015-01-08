Feature: git kill: don't remove a given non-feature branch (with open changes)

  (see ../current_branch/non_feature_branch_with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And non-feature branch configuration "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill qa` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | qa      | local and remote | qa commit   | qa_file   |
      | feature | local and remote | good commit | good_file |
