Feature: git kill: does not remove non-feature branches (without open changes)

  As a developer accidentally running "git kill" on a non-feature branch
  I want the command to not perform the operation
  So that my release infrastructure remains intact and my project remains shippable.


  Background:
    Given I have a feature branch named "good-feature"
    Given non-feature branch configuration "qa"
    And the following commits exist in my repository
      | branch       | location         | message     | file name |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
    And I am on the "qa" branch
    When I run `git kill` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "qa" branch
    And the existing branches are
      | repository | branches               |
      | local      | main, qa, good-feature |
      | remote     | main, qa, good-feature |
    And I have the following commits
      | branch       | location         | message     | files     |
      | good-feature | local and remote | good commit | good_file |
      | qa           | local and remote | qa commit   | qa_file   |
