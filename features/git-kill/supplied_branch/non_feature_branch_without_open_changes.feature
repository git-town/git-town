Feature: git kill: don't remove a given non-feature branch (without open changes)

  Background:
    Given I have a feature branch named "good-feature"
    Given non-feature branch configuration "qa"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
    And I am on the "good-feature" branch
    When I run `git kill qa` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "good-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES               |
      | local      | main, qa, good-feature |
      | remote     | main, qa, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | qa           | local and remote | qa commit   | qa_file   |
      | good-feature | local and remote | good commit | good_file |
