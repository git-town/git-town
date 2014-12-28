Feature: git kill: don't remove the main branch (with open changes)

  As a developer accidentally trying to kill the main branch
  I should see an error that I cannot delete the main branch
  So that my main development line remains intact and my project stays shippable.


  Background:
    Given I have a feature branch named "good-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE     | FILE NAME |
      | good-feature | local and remote | good commit | good_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |
