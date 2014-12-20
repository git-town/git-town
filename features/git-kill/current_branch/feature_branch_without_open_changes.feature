Feature: Git Kill: Killing the current feature branch without open changes

  Background:
    Given I have feature branches named "good-feature" and "dead-feature"
    And the following commits exist in my repository
      | branch       | location         | message         | file name        |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "dead-feature" branch
    When I run `git kill`


  Scenario: result
    Then I end up on the "main" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: Undoing the kill
    When I run `git kill --undo`
    Then I end up on the "dead-feature" branch
    And the existing branches are
      | repository | branches                         |
      | local      | main, dead-feature, good-feature |
      | remote     | main, dead-feature, good-feature |
    And I have the following commits
      | branch       | location         | message         | files            |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
