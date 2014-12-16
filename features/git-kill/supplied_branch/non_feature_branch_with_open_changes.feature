Feature: git kill: don't remove non-feature branches (with open changes)

  As a developer accidentally trying to kill a non-feature branch
  I should be warned that this branch type can not be removed
  So that my release infrastructure remains intact.


  Background:
    Given I have a feature branch named "good-feature"
    Given non-feature branch configuration "qa"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
    And I am on the "good-feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill qa` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "good-feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES               |
      | local      | main, qa, good-feature |
      | remote     | main, qa, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES   |
      | qa           | local and remote | qa commit   | qa_file |
      | good-feature | local and remote | good commit | good_file |
