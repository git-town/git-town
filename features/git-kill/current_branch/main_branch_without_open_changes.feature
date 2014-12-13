Feature: git kill: don't remove the main branch (without open changes)

  As a developer accidentally trying to kill the main branch
  I should be warned that the command does not remove the main branch
  So that my main development line remains intact and my project stays shippable.


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
