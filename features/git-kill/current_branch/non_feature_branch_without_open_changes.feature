Feature: git kill: don't remove non-feature branches (without open changes)

  Background:
    Given I have a feature branch named "feature"
    And non-feature branch configuration "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
    And I am on the "qa" branch
    When I run `git kill` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "qa" branch
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILES     |
      | feature | local and remote | good commit | good_file |
      | qa      | local and remote | qa commit   | qa_file   |
