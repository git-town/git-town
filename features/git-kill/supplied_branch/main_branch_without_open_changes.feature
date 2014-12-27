Feature: git kill: don't remove the main branch (with open changes)

  Background:
    Given I have a feature branch named "good-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
      | main         | local and remote | main commit | main_file |
    And I am on the "good-feature" branch
    When I run `git kill main` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "good-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |
      | main         | local and remote | main commit | main_file |
