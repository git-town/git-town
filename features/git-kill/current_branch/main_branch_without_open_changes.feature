Feature: git kill: does not remove the main branch (without open changes)

  As a developer accidentally running "git kill" on the main branch
  I want the command to not perform the operation
  So that my main development line remains intact and my project remains shippable.


  Background:
    Given I have a feature branch named "good-feature"
    And I am on the "main" branch
    When I run `git kill` while allowing errors


  Scenario: result
    Then I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And the existing branches are
      | repository | branches           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
