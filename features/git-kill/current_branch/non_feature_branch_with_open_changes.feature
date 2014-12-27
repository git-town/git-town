Feature: git kill: don't remove non-feature branches (with open changes)

  Background:
    Given I have a feature branch named "good-feature"
    Given non-feature branch configuration "qa"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
    And I am on the "qa" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES               |
      | local      | main, qa, good-feature |
      | remote     | main, qa, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
