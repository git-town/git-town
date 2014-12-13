Feature: git kill: does not remove the main branch (with open changes)

  As a developer accidentally running "git kill" while developing on the main branch
  I want the command to not perform the operation
  So that my main development line remains intact and my project remains shippable.


  Background:
    Given I have a feature branch named "good-feature"
    And the following commits exist in my repository
      | branch       | location         | message     | file name |
      | good-feature | local and remote | good commit | good_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git kill` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
